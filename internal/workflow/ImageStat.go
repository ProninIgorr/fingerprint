package workflow

import (
	"github.com/ProninIgorr/fingerprint/internal/filestat"
	"github.com/ProninIgorr/fingerprint/internal/imgs/matrix"
)

type ImageStat struct {
	Fs filestat.FileStat
	IM *matrix.M
}
