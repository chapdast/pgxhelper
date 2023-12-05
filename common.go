package pgxhelper

type Dir string

var (
	DirAscending  Dir = "ASC"
	DirDescending Dir = "DESC"
)

type Sort struct {
	Column    string
	Direction Dir
}
