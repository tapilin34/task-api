package storage

import (
	"errors"
	"sort"
	"sync"
	"tasks-api/internal/models"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskCreate = errors.New("task not create")

type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
type Cache struct {
	tasks  map[int]models.Task
	mu     sync.RWMutex
	nextID int
}

func NewCache() *Cache {
	return &Cache{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (c *Cache) List() []models.Task {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// как сортировать map подсказал ИИ
	ids := make([]int, 0, len(c.tasks))
	for id := range c.tasks {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	result := make([]models.Task, 0, len(ids))
	for _, id := range ids {
		result = append(result, c.tasks[id])
	}

	return result
}

func (c *Cache) Create(task models.Task) (models.Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	task.ID = c.nextID
	c.nextID++
	if task.CreatedAt == "" {
		task.CreatedAt = time.Now().Format(time.RFC3339)
	}
	c.tasks[task.ID] = task
	if _, ok := c.tasks[task.ID]; !ok {
		// сюда мы попадём только если что‑то пошло совсем не так
		return models.Task{}, ErrTaskCreate
	}
	return task, nil
}

func (c *Cache) Get(id int) (models.Task, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	task, ok := c.tasks[id]
	return task, ok
}

func (c *Cache) Update(id int, task models.Task) (models.Task, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	oldTask, ok := c.tasks[id]
	if !ok {
		return models.Task{}, ErrTaskNotFound
	}
	task.ID = oldTask.ID
	task.CreatedAt = oldTask.CreatedAt
	if task.CreatedAt == "" {
		task.CreatedAt = time.Now().Format(time.RFC3339)
	}
	c.tasks[id] = task
	return task, nil
}

func (c *Cache) Delete(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}
	delete(c.tasks, id)
	return nil
}
