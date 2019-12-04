package protocol

type CHandler interface {
	CListHandle(data []byte)

	CListUppageHandle(data []byte)

	CUploadbodyNextpackHandle(data []byte)

	CDownloadheadHandle(data []byte)
	CDownloadbodyHandle(data []byte)

	CListdisk(data []byte)
}

type SHandler interface {
	SListHandle(data []byte)

	SListUppageHandle(data []byte)

	SUploadheadHandle(data []byte)
	SUploadbodyHandle(data []byte)

	SDownloadheadHandle(data []byte)
	SDownloadbodyHandle(data []byte)

	SDownloadbodyNextpackHandle(data []byte)

	SListdisk(data []byte)
}
