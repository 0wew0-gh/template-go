package tool

import "strings"

type StrMenu struct {
	ID       string
	Name     string
	Children []StrMenu
}

var (
	ClassPathSeparador string = "#@#"
)

// 递归生成菜单
func insertPath(root *StrMenu, parts []string) {
	if len(parts) == 0 {
		return
	}

	current := parts[0]

	// 查找是否已有该节点
	for i := range root.Children {
		if root.Children[i].Name == current {
			insertPath(&root.Children[i], parts[1:])
			return
		}
	}

	// 没有则新建
	newNode := StrMenu{Name: current}
	root.Children = append(root.Children, newNode)
	insertPath(&root.Children[len(root.Children)-1], parts[1:])
}

// 把路径列表生成树结构
// 例：/a/b/c/d
// StrMenu{
// 	Name: "ROOT",
// 	Children: []StrMenu{
// 		Name: "a",
// 		Children: []StrMenu{
// 			{
// 				Name: "b",
// 				Children: []StrMenu{
// 					{
// 						Name: "c", Children: []StrMenu{
// 							{
// 								Name: "d",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }
func BuildTree(paths []string) StrMenu {
	root := StrMenu{Name: "ROOT"}
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		parts := strings.Split(p, ClassPathSeparador)
		insertPath(&root, parts)
	}
	return root
}

// 递归插入 path，对最里层节点设置 leafID，其它节点 ID 默认 "0"
// 递归插入节点
func insertPathWithID(root *StrMenu, parts []string, id string) {
	if len(parts) == 0 {
		return
	}

	current := parts[0]

	// 查找是否已有子节点
	var child *StrMenu
	for i := range root.Children {
		if root.Children[i].Name == current {
			child = &root.Children[i]
			break
		}
	}

	// 若不存在则创建
	if child == nil {
		newNode := StrMenu{
			Name: current,
			ID:   "0",
		}
		root.Children = append(root.Children, newNode)
		child = &root.Children[len(root.Children)-1]
	}

	// 最深层赋值 ID
	if len(parts) == 1 {
		child.ID = id
		return
	}

	insertPathWithID(child, parts[1:], id)
}

// 把路径列表生成树结构，并设置 ID
// 例：/a/b/c/d []string{1}
// StrMenu{
// 	Name: "ROOT",
// 	Children: []StrMenu{
//		ID: 0,
// 		Name: "a",
// 		Children: []StrMenu{
// 			{
//				ID: 0,
// 				Name: "b",
// 				Children: []StrMenu{
// 					{
//						ID: 0,
// 						Name: "c", Children: []StrMenu{
// 							{
//								ID: 1,
// 								Name: "d",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }
func BuildTreeWithIDs(paths, ids []string) StrMenu {
	root := StrMenu{Name: "root", ID: "0"}

	for i, p := range paths {
		parts := strings.Split(p, ClassPathSeparador)
		insertPathWithID(&root, parts, ids[i])
	}

	return root
}
