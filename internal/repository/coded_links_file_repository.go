package repository

import (
	"bufio"
	"errors"
	"os"
	"sync"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/model"
)

type CodedLinksFileRepository struct {
	idsMap          map[int]*model.CodedLink
	shortUrlsMap    map[string]*model.CodedLink
	originalUrlsMap map[string]*model.CodedLink
	nextID          int
	file            *os.File
	writer          *bufio.Writer
	reader          *bufio.Reader
	rwMutex         *sync.RWMutex
}

func NewCodedLinksFileRepository(cfg *config.Config) (*CodedLinksFileRepository, error) {
	file, err := os.OpenFile(*cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.New("не удалось открыть файл для хранения ссылок по пути `" + *cfg.FileStoragePath + "`: " + err.Error())
	}

	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(file)

	repository := &CodedLinksFileRepository{
		idsMap:          make(map[int]*model.CodedLink),
		shortUrlsMap:    make(map[string]*model.CodedLink),
		originalUrlsMap: make(map[string]*model.CodedLink),
		file:            file,
		writer:          writer,
		reader:          reader,
		nextID:          1_000_000,
		rwMutex:         &sync.RWMutex{},
	}

	err = repository.loadFromFile()
	if err != nil {
		return nil, errors.New("не удалось загрузить данные из файла: " + err.Error())
	}

	return repository, nil
}

func (lr *CodedLinksFileRepository) loadFromFile() error {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	for {
		line, err := lr.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, os.ErrClosed) || errors.Is(err, os.ErrInvalid) {
				return errors.New("не удалось прочитать данные из файла: " + err.Error())
			}
			break
		}

		codedLink, err := model.FromJSON(line)
		if err != nil {
			return errors.New("не удалось распарсить строку из файла в структуру CodedLink: " + err.Error())
		}

		lr.idsMap[codedLink.UUID] = codedLink
		lr.shortUrlsMap[codedLink.ShortURL] = codedLink
		lr.originalUrlsMap[codedLink.OriginalURL] = codedLink

		if codedLink.UUID >= lr.nextID {
			lr.nextID = codedLink.UUID + 1
		}
	}

	return nil
}

func (lr *CodedLinksFileRepository) Close() error {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	if err := lr.writer.Flush(); err != nil {
		return errors.New("не удалось записать данные в файл: " + err.Error())
	}

	if err := lr.file.Close(); err != nil {
		return errors.New("не удалось закрыть файл: " + err.Error())
	}

	return nil
}

func (lr *CodedLinksFileRepository) CreateCodedLink(codedLink *model.CodedLink) error {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	if _, isOk := lr.shortUrlsMap[codedLink.OriginalURL]; isOk {
		return nil
	}

	lr.idsMap[codedLink.UUID] = codedLink
	lr.shortUrlsMap[codedLink.ShortURL] = codedLink
	lr.originalUrlsMap[codedLink.OriginalURL] = codedLink
	line, err := codedLink.ToJSON()
	if err != nil {
		return errors.New("не удалось сериализовать структуру CodedLink в строку для записи в файл: " + err.Error())
	}

	_, err = lr.writer.Write(line)
	if err != nil {
		return err
	}
	_, err = lr.writer.WriteRune('\n')
	if err != nil {
		return errors.New("не удалось сериализовать структуру CodedLink в строку для записи в файл: " + err.Error())
	}
	err = lr.writer.Flush()
	if err != nil {
		return errors.New("не удалось сериализовать структуру CodedLink в строку для записи в файл: " + err.Error())
	}

	return nil
}

func (lr *CodedLinksFileRepository) GetByID(id int) (link *model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.idsMap[id]
	return result, isOk
}

func (lr *CodedLinksFileRepository) GetByShortURL(shortLink string) (link *model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.shortUrlsMap[shortLink]
	return result, isOk
}

func (lr *CodedLinksFileRepository) GetByOriginalURL(originalLink string) (link *model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.originalUrlsMap[originalLink]
	return result, isOk
}

func (lr *CodedLinksFileRepository) GetUniqueID() int {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	result := lr.nextID
	lr.nextID++
	return result
}
