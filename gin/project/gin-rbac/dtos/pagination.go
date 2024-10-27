package dtos

type Filter struct {
	Field string `json:"field"`                                                                  // 过滤字段
	Value string `json:"value"`                                                                  // 过滤值
	Op    string `form:"op" binding:"required,oneof=eq neq gt gte lt lte contains not_contains"` // 操作符
}

type PaginationReqDTO struct {
	Page    int      `json:"page" form:"page" binding:"omitempty"`       // 当前页码
	Size    int      `json:"size" form:"size" binding:"omitempty"`       // 每页数量
	SortBy  string   `json:"sort_by" form:"sort_by" binding:"omitempty"` // 排序字段
	Order   string   `json:"order" form:"order" binding:"omitempty"`     // 排序顺序
	Filters []Filter `json:"filters" form:"filters" binding:"omitempty"` // 过滤条件
}

// GetDefaultPaginationDTO 获取默认分页参数
func GetDefaultPaginationDTO(paginationDTO PaginationReqDTO) PaginationReqDTO {
	defaultPage := 1       // 默认页码
	defaultSize := 10      // 默认每页数量
	defaultOrder := "desc" // 默认排序顺序
	defaultSortBy := "id"  // 默认排序字段

	// 如果 paginationDTO.Page 是零值且不符合约束（> 0），则使用默认值
	if paginationDTO.Page <= 0 {
		paginationDTO.Page = defaultPage
	}

	// 如果 paginationDTO.PageSize 是零值或者不符合约束（> 0 和 <= 100），则使用默认值
	if paginationDTO.Size <= 0 || paginationDTO.Size > 100 {
		paginationDTO.Size = defaultSize
	}

	// 如果 paginationDTO.Order 为空，则使用默认值
	if paginationDTO.Order == "" {
		paginationDTO.Order = defaultOrder
	}

	// 如果 paginationDTO.SortBy 为空，则使用默认值
	if paginationDTO.SortBy == "" {
		paginationDTO.SortBy = defaultSortBy
	}

	// 根据转换后的值或默认值创建并返回 PaginationDTO 实例
	return paginationDTO
}

type PaginationResult[T any] struct {
	Total int64 `json:"total"` // 总记录数
	Items []T   `json:"items"` // 当前页数据项
	Page  int   `json:"page"`  // 当前页码
	Size  int   `json:"size"`  // 每页数量
	Pages int64 `json:"pages"` // 总页数
}

// PaginationResultExample 分页结果示例，由于swagger不支持泛型，所以使用示例结构体代替
type PaginationResultExample struct {
	Total int64         `json:"total"`
	Items []interface{} `json:"items"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Pages int64         `json:"pages"`
}

// NewPaginationResult 创建分页结果实例
func NewPaginationResult[T any](total int64, data []T, page, size int) PaginationResult[T] {
	totalPages := (total + int64(size) - 1) / int64(size)
	return PaginationResult[T]{
		Total: total,
		Items: data,
		Page:  page,
		Size:  size,
		Pages: totalPages,
	}
}
