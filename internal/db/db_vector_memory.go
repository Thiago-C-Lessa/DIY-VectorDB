package db

import (
	"DIY-VectorDB/internal/exceptions"
	"DIY-VectorDB/internal/http/embedding"
	"DIY-VectorDB/internal/models"
	"DIY-VectorDB/internal/utils"
	"sort"

	"encoding/json"
	"math"
	"math/rand/v2"
	"sync"
)

var (
	M uint64 = 8 //para um dataset <= 10k
)

type Node struct {
	level     uint64
	embedding [768]float32
	data      *json.RawMessage
	key       *string
	neighbors [][]*Node
}

type VecMemDB struct {
	sync.RWMutex
	maxLevel           uint64
	maxNeighbors       uint64
	eFConstruction     uint64
	eFSearch           uint64
	entry              *Node
	embeddings         map[string][768]float32
	dataForSingleQuery map[[768]float32]*Node
}

type candidate struct {
	cosProduct float32
	node       *Node
}

func NewVecMemDB() *VecMemDB {
	return &VecMemDB{
		embeddings:         make(map[string][768]float32),
		dataForSingleQuery: make(map[[768]float32]*Node),
		maxNeighbors:       M,
		eFConstruction:     2 * M,
		eFSearch:           M * M,
		entry:              nil,
		maxLevel:           0,
	}
}

func newNode(embedding [768]float32, data *json.RawMessage, key *string) *Node {
	n := new(Node)

	n.embedding = embedding
	n.level = randomLevel(float64(M))
	n.data = data
	n.neighbors = make([][]*Node, n.level+1)
	n.key = key

	return n
}

func (n *Node) GetEmbedding() [768]float32 {
	return n.embedding
}

func randomLevel(m float64) uint64 {
	var maxLevel uint64 = 64
	f := rand.Float64()
	level := uint64(-math.Floor(math.Log(f) * m))

	if level > maxLevel {
		level = maxLevel
	}

	return level
}

func (vdb *VecMemDB) Insert(key string, value json.RawMessage) error {
	vdb.Lock()
	defer vdb.Unlock()

	_, ok := vdb.embeddings[key]
	if ok {
		return &exceptions.ErrorKeyAlreadyExists{Key: key}
	}

	emb, err := embedding.GenereteEmbedding(key)
	if err != nil {
		return err
	}

	vdb.embeddings[key] = emb
	n := newNode(emb, &value, &key)
	vdb.dataForSingleQuery[emb] = n

	if vdb.entry == nil {
		vdb.entry = n
		vdb.maxLevel = n.level
		return nil
	}

	//determina o entry point
	ep := vdb.entry
	for i := int(vdb.maxLevel); i >= int(n.level); i-- {
		ep = vdb.greedSearch(n.embedding, ep, uint64(i))
	}

	//conecta os nós vizinhos
	for i := int(utils.MinUint64(n.level, vdb.maxLevel)); i >= 0; i-- {
		candiates := vdb.searchLayer(n.embedding, ep, uint64(i), vdb.eFConstruction) //eFSearch é para busca
		neighbors := vdb.selectNeighbors(candiates, n.embedding, vdb.maxNeighbors)

		n.neighbors[i] = neighbors

		for _, nb := range neighbors {
			nb.neighbors[i] = append(nb.neighbors[i], n)

			if len(nb.neighbors[i]) > int(vdb.maxNeighbors) {
				nb.neighbors[i] = vdb.selectNeighbors(nb.neighbors[i], nb.embedding, vdb.maxNeighbors)
			}
		}
	}

	if n.level > vdb.maxLevel {
		vdb.maxLevel = n.level
		vdb.entry = n
	}

	return nil
}

func (vdb *VecMemDB) greedSearch(target [768]float32, entry *Node, level uint64) *Node {
	closestNode := entry
	var currentDist float32 = utils.CosineProductPreNormalized(target, entry.embedding)
	var dist float32

	changed := true
	for changed {
		changed = false

		for _, n := range closestNode.neighbors[level] {
			dist = utils.CosineProductPreNormalized(target, n.embedding)

			if dist < currentDist {
				continue
			} else {
				currentDist = dist
				closestNode = n
				changed = true
			}
		}
	}

	return closestNode
}

func (vdb *VecMemDB) searchLayer(target [768]float32, entry *Node, level, radius uint64) []*Node {

	visited := make(map[*Node]bool)
	candidates := []*Node{entry}
	results := []*Node{entry}

	visited[entry] = true

	for len(candidates) > 0 && uint64(len(results)) < radius {
		curr := candidates[0]
		candidates = candidates[1:]

		for _, n := range curr.neighbors[level] {
			if visited[n] {
				continue
			}
			visited[n] = true

			candidates = append(candidates, n)
			results = append(results, n)
		}
	}

	// ordena por similaridade
	results = SortBySimilarity(results, target)

	if uint64(len(results)) > radius {
		results = results[:radius]
	}

	return results
}

func (vdb *VecMemDB) selectNeighbors(nodes []*Node, target [768]float32, max uint64) []*Node {

	SortBySimilarity(nodes, target)

	if uint64(len(nodes)) > max {
		nodes = nodes[:max]
	}

	return nodes
}

func SortBySimilarity(v []*Node, target [768]float32) []*Node {
	result := make([]*Node, len(v))
	c := make([]candidate, len(v))

	for i, _ := range v {
		c[i] = candidate{
			cosProduct: utils.CosineProductPreNormalized(v[i].GetEmbedding(), target),
			node:       v[i],
		}
	}

	sort.Slice(c, func(i, j int) bool {
		return c[i].cosProduct > c[j].cosProduct
	})

	for i, _ := range v {
		result[i] = c[i].node
	}

	return result
}

func (vdb *VecMemDB) SelectSimilar(key string, k uint64) (models.ResponseData, error) {
	vdb.RLock()
	defer vdb.RUnlock()

	if k == 0 {
		k = 8
	}

	var response models.ResponseData

	query, ok := vdb.embeddings[key]
	if !ok {
		return models.ResponseData{}, &exceptions.ErrorContentNotFound{Key: key}
	}

	if vdb.entry == nil {
		return response, nil
	}

	ep := vdb.entry

	for level := int(vdb.maxLevel); level > 0; level-- {
		ep = vdb.greedSearch(query, ep, uint64(level))
	}

	candidates := vdb.searchLayer(query, ep, 0, vdb.eFSearch)

	results := vdb.selectNeighbors(candidates, query, k)

	for _, n := range results {
		response.Keys = append(response.Keys, *n.key)
		response.Contents = append(response.Contents, *n.data)

	}

	return response, nil
}

func (vdb *VecMemDB) Select(key string) (models.ResponseData, error) {
	vdb.RLock()
	defer vdb.RUnlock()

	var response models.ResponseData

	emb, ok := vdb.embeddings[key]
	if !ok {
		return response, &exceptions.ErrorContentNotFound{Key: key}
	}

	d, _ := vdb.dataForSingleQuery[emb]

	response.Keys = append(response.Keys, key)
	response.Contents = append(response.Contents, *d.data)

	return response, nil
}

func (vdb *VecMemDB) ListAll() (models.ResponseData, error) {
	vdb.RLock()
	defer vdb.RUnlock()

	var data models.ResponseData

	for key, emb := range vdb.embeddings {
		content, ok := vdb.dataForSingleQuery[emb]
		if !ok {
			return data, &exceptions.ErrorContentNotFound{Key: key}
		}

		data.Keys = append(data.Keys, key)
		data.Contents = append(data.Contents, *content.data)
	}

	return data, nil
}

func (vdb *VecMemDB) Update(key string, value json.RawMessage) error {
	vdb.Lock()
	defer vdb.Unlock()
	emb, ok := vdb.embeddings[key]
	if !ok {
		return &exceptions.ErrorContentNotFound{Key: key}
	}

	vdb.dataForSingleQuery[emb].data = &value
	return nil
}
