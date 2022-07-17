package services

func UpdateMeta(files *[]File) {
	storage := GetStorage()
	fc, ss := storage.GetMeta()
	fc += int64(len(*files))
	for _, f := range *files {
		ct := f.Size - f.NSize
		if ct > 0 {
			ss += ct
		}
	}
	storage.SetMeta(fc, ss)
}
