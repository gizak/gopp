package plugin
import p0 "github.com/gizak/gopp/plugin/parsepanic"
import p1 "github.com/gizak/gopp/plugin/gocolor"
func init() {
 	Plugins = append(Plugins,PluginEntry{Name:"parsepanic",Plugin:p0.PLUGIN,Uri:"github.com/gizak/gopp/plugin/parsepanic"})
	Plugins = append(Plugins,PluginEntry{Name:"gocolor",Plugin:p1.PLUGIN,Uri:"github.com/gizak/gopp/plugin/gocolor"})
}