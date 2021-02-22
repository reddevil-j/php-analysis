// Package visitor contains walker.visitor implementations
package visitor

import (
	"reflect"
	"io"
	"strings"
	"bytes"
	"github.com/z7zmey/php-parser/parser"
	"github.com/z7zmey/php-parser/printer"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/name"
	"github.com/z7zmey/php-parser/node/stmt"
	l "medium/analyze/logger"
	"github.com/z7zmey/php-parser/walker"
)



// Dumper writes ast hierarchy to an io.Writer
// Also prints comments and positions attached to nodes


var cname string = "" // hold the name of class
var ename string = "" // hold the name of extended class
var functionname string
var methodname string

type DefWalker struct {
	Writer io.Writer
	Indent string
	Comments parser.Comments
	Positions parser.Positions
	NsResolver *NamespaceResolver
}

func NodeSource(n *node.Node) string {
	out := new(bytes.Buffer)
	p := printer.NewPrinter(out, "    ")
	p.Print(*n)
	return strings.Replace(out.String(), "\n", "\\ ", -1)

}

var File string
var RelativePath string
var Functions = make(map[string]bool)
var Methods = make(map[string]bool)
// EnterNode is invoked at every node in hierarchy
func (d DefWalker) EnterNode(w walker.Walkable) bool {

	n := w.(node.Node)
	fname := RelativePath

	switch reflect.TypeOf(n).String() {

	case "*stmt.Class":
		class := n.(*stmt.Class)
		if namespacedName, ok := d.NsResolver.ResolvedNames[class]; ok {
			cname = namespacedName
		} else {
			className, ok := class.ClassName.(*node.Identifier)
			if !ok {
				l.Log(l.Error, "class name is not resolved: %s:%s", NodeSource(&n), File)
				break
			}
			cname = className.Value
		}
		extClass := class.Extends
		if extClass != nil {
			if namespacedName, ok := d.NsResolver.ResolvedNames[extClass]; ok {
				ename = namespacedName
			} else {
				extName, ok := extClass.ClassName.(*name.Name)
				if !ok {
					l.Log(l.Error, "extended class name is not resolved: %s", NodeSource(&n), File)
					break
				}
				ename = extName.Parts[len(extName.Parts)-1].(*name.NamePart).Value
			}
		}
		break
	case "*stmt.Function":
		function := n.(*stmt.Function)
		funcName, ok := function.FunctionName.(*node.Identifier)
		if ok {
			if namespacedName, ok := d.NsResolver.ResolvedNames[function]; ok {
				functionname = namespacedName
				Functions[functionname +"|"+fname] = true
				_ = functionname
			} else {
				functionname = funcName.Value
				Functions[functionname+"|"+fname] = true
				_ = functionname
			}
		} else {
			l.Log(l.Error, "function name is not resolved: %s:%s", NodeSource(&n), File)
		}
		break
	case "*stmt.ClassMethod":
		classmethod := n.(*stmt.ClassMethod)
		mname, ok := classmethod.MethodName.(*node.Identifier)
		if ok {
			if namespacedName, ok := d.NsResolver.ResolvedNames[classmethod]; ok {
				methodname = namespacedName
				Methods[namespacedName+"|"+fname] = true
			} else {
				methodname = cname + "\\" + mname.Value
				_ = methodname
				Methods[methodname+"|"+fname] = true
			}
		} else {
			l.Log(l.Error, "method name is not resolved: %s:%s", NodeSource(&n), File)
		}
		break
	}
	return true
}


// GetChildrenVisitor is invoked at every node parameter that contains children nodes
func (d DefWalker) GetChildrenVisitor(key string) walker.Visitor {
	return DefWalker{d.Writer, d.Indent + "    ", d.Comments, d.Positions, d.NsResolver }
}

// LeaveNode is invoked after node process
func (d DefWalker) LeaveNode(w walker.Walkable) {
	//parse := false
	n := w.(node.Node)

	switch reflect.TypeOf(n).String() {
	case "*stmt.Class":
		cname = ""
		break

	}
}

