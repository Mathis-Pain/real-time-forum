package comments

import "fmt"

func IsValidFormat(comment Comment) error {
	if comment.Content == "" {
		return fmt.Errorf("no content provided")
	}

	return nil
}
