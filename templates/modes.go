package templates

import "variant/log"

type LocalCtx struct {
	*BaseCtx
	KeyVar   string // Key 变量名称
	KeyValue string // Key 值
	IVVar    string // IV 变量名称
	IVValue  string // IV 值
	Payload  string // Payload 加密后的值
}

// SetLocal 设置 LocalCtx 并绑定 BaseCtx
func (b *TmplBuilder) SetLocal(base *BaseCtx, fn func(*LocalCtx)) *TmplBuilder {
	// 如果 Local 为 nil，主动创建新实例
	if b.tmpl.Local == nil {
		b.tmpl.Local = &LocalCtx{}
	}

	// 绑定 BaseCtx（即使 Base 为 nil 也接受）
	b.tmpl.Local.BaseCtx = base

	// 应用用户配置函数（如果提供）
	if fn != nil {
		fn(b.tmpl.Local)
	}

	return b
}

// DynamicCtx 包含动态加密相关的上下文信息
type DynamicCtx struct {
	*BaseCtx // 基础上下文

	KeyVar   string // 密钥变量名，如"aesKey"
	KeyValue string // 密钥值
	IVVar    string // 初始化向量变量名，如"aesIV"
	IVValue  string // 向量值
	Payload  string // 加密后的Payload值
}

// SetDynamic 设置动态加密上下文
// 参数:
//
//	base: 基础上下文，可为nil
//	fn: 配置函数，用于自定义DynamicCtx，可为nil
//
// 返回值:
//
//	*TmplBuilder: 返回构建器实例以支持链式调用
func (b *TmplBuilder) SetDynamic(base *BaseCtx, fn func(*DynamicCtx)) *TmplBuilder {
	// 确保Dynamic上下文已初始化
	if b.tmpl.Dynamic == nil {
		b.tmpl.Dynamic = &DynamicCtx{}
	}

	// 设置基础上下文
	b.tmpl.Dynamic.BaseCtx = base

	// 应用用户自定义配置
	if fn != nil {
		fn(b.tmpl.Dynamic)
	}

	return b
}

type ArgsCtx struct {
	*BaseCtx
	Import string
	Passwd string
}

func (b *TmplBuilder) SetArgs(base *BaseCtx, fn func(*ArgsCtx)) *TmplBuilder {
	if b.tmpl.Args == nil {
		b.tmpl.Args = &ArgsCtx{}
	}

	b.tmpl.Args.BaseCtx = base
	if fn != nil {
		fn(b.tmpl.Args)
	}

	return b
}

type PokemonCtx struct {
	*BaseCtx
	MainPokemon    string
	PokemonPayload []string
	DecryptPokemon string
}

// lazySetContext 是通用方法，适用于所有子上下文类型。
// T 必须是 *ArgsCtx、*LocalCtx、*DynamicCtx、*PokemonCtx 之一。
func lazySetContext[T any](
	b *TmplBuilder,
	target **T, // 指向字段的指针（如 &b.tmpl.Args）
	base *BaseCtx,
	fn func(*T),
) *TmplBuilder {
	if *target == nil {
		*target = new(T)
	}

	// 所有子上下文都嵌入了 BaseCtx
	switch v := any(*target).(type) {
	case *ArgsCtx:
		v.BaseCtx = base
	case *LocalCtx:
		v.BaseCtx = base
	case *DynamicCtx:
		v.BaseCtx = base
	case *PokemonCtx:
		v.BaseCtx = base
	default:
		log.Fatalf("unsupported context type: %T", *target)
	}

	if fn != nil {
		fn(*target)
	}

	return b
}
