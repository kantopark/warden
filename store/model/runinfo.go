package model

// The RunInfo contains information on the how to run the function.
// Specifically, it links the alias to the commit. The default alias
// is the "latest" tag, thus, it would mean a request to the base
// url of the function will hit that specific runtime. Internally,
// the default empty string "" will be aliased to "latest" as well.
type RunInfo struct {
	Id         uint   `gorm:"primary_key"`
	Alias      string `gorm:"unique_index:idx_alias_function"`
	Commit     string `gorm:"varchar(100)"`
	FunctionID uint   `gorm:"unique_index:idx_alias_function"`
}
