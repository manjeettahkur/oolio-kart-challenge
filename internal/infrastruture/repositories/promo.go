package repositories

import (
	"bufio"
	"context"
	"fmt"
	"ooliokartchallenge/internal/domain/interfaces"
	"os"
	"strings"
	"sync"
)

type PromoRepository struct {
	filePaths []string
	mutex     sync.RWMutex
}

func NewPromoRepository(filePaths []string) interfaces.PromoRepository {
	repo := &PromoRepository{
		filePaths: filePaths,
	}

	fmt.Printf("Initializing promo repository with %d files...\n", len(filePaths))
	for i, path := range filePaths {
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("Warning: File %d (%s) not accessible: %v\n", i+1, path, err)
		} else {
			fmt.Printf("File %d: %s - ready\n", i+1, path)
		}
	}
	fmt.Printf("Promo repository initialized (files will be scanned on-demand)\n")

	return repo
}

func (r *PromoRepository) ValidateCode(ctx context.Context, code string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.filePaths) < 2 {
		return false, fmt.Errorf("at least 2 coupon files required for validation")
	}


	type scanResult struct {
		fileIndex int
		found     bool
		err       error
	}

	resultChan := make(chan scanResult, len(r.filePaths))

	for i, filePath := range r.filePaths {
		go func(index int, path string) {
			found, err := r.scanFileForCode(ctx, path, code)
			resultChan <- scanResult{
				fileIndex: index,
				found:     found,
				err:       err,
			}
		}(i, filePath)
	}

	filesWithCode := 0
	for i := 0; i < len(r.filePaths); i++ {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case result := <-resultChan:
			if result.err != nil {
				fmt.Printf("Error scanning file %d: %v\n", result.fileIndex+1, result.err)
				continue
			}
			if result.found {
				filesWithCode++

				if filesWithCode >= 2 {
					return true, nil
				}
			}
		}
	}

	return filesWithCode >= 2, nil
}

func (r *PromoRepository) scanFileForCode(ctx context.Context, filePath, targetCode string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineCount := 0
	for scanner.Scan() {

		if lineCount%10000 == 0 {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			default:
			}
		}

		code := strings.TrimSpace(scanner.Text())
		if code == targetCode {
			return true, nil
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return false, nil
}
