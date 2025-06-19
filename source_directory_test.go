package iterate

import (
       "context"
       "testing"
       "path/filepath"
       "fmt"
)

func TestDirectorySource(t *testing.T) {

     ctx := context.Background()

     abs_path, err := filepath.Abs("fixtures")

     if err != nil {
     	t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
     }

     src, err := NewDirectorySource(ctx, "directory://")

     if err != nil {
     	t.Fatalf("Failed to create new directory source, %v", err)
     }

     for rec, err := range src.Walk(ctx, abs_path) {

     	 if err != nil {
	    //t.Fatalf("Failed to walk '%s', %v", abs_path, err)
	    //break
	}
	    
     	fmt.Println(rec.Path)
     }
}