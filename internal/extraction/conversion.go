package extraction

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func PDFtoJPEGinMem(ctx context.Context, coverPDFBytes []byte) ([]byte , error) {

	cmd := exec.CommandContext(ctx, "pdftoppm", "-jpeg", "-singlefile" , "-r", "150", "-", "-")
	
	cmd.Stdin = bytes.NewReader(coverPDFBytes)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF to JPEG: %w\nstderr: %s", err, stderrBuf.String())
	}

	return stdoutBuf.Bytes(), nil

}

