// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/melihkorkmaz/gtd/internal/views/partials"

func Base(title string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\" data-theme=\"light\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/layouts/base.templ`, Line: 11, Col: 16}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</title><!-- DaisyUI with Tailwind CSS --><link href=\"https://cdn.jsdelivr.net/npm/daisyui@3.9.4/dist/full.css\" rel=\"stylesheet\" type=\"text/css\"><script src=\"https://cdn.tailwindcss.com\"></script><!-- Alpine.js --><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.13.3/dist/cdn.min.js\"></script><!-- HTMX --><script src=\"https://unpkg.com/htmx.org@1.9.6\"></script><!-- Custom CSS --><link rel=\"stylesheet\" href=\"/static/css/main.css\"></head><body class=\"min-h-screen bg-base-200\"><div class=\"container mx-auto p-4\"><header class=\"navbar bg-base-100 rounded-box shadow-lg mb-6\"><div class=\"flex-1\"><a href=\"/\" class=\"btn btn-ghost text-xl\">GTD App</a></div><div class=\"flex-none\"><!-- Search box --><div class=\"form-control mx-2\"><form action=\"/tasks/search\" method=\"GET\" class=\"flex\" hx-get=\"/tasks/search\" hx-trigger=\"submit\" hx-target=\"#main-content\" hx-swap=\"innerHTML\"><div class=\"relative\"><input type=\"text\" name=\"q\" placeholder=\"Search tasks...\" class=\"input input-bordered w-24 md:w-auto\" hx-get=\"/tasks/search\" hx-trigger=\"keyup changed delay:500ms\" hx-target=\"#search-results\" hx-indicator=\".search-indicator\"><div class=\"search-indicator\"><span class=\"loading loading-spinner loading-xs\"></span></div></div><button type=\"submit\" class=\"btn btn-primary\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z\"></path></svg></button></form></div><!-- Quick Capture Button --><button class=\"btn btn-success btn-sm mx-2\" onclick=\"document.getElementById(&#39;quick-capture-modal&#39;).showModal()\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5 mr-1\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Capture</button><ul class=\"menu menu-horizontal px-1\"><li><a href=\"/tasks\">All Tasks</a></li><li><a href=\"/tasks?status=inbox\">Inbox</a></li><li><a href=\"/tasks?status=next\">Next Actions</a></li><li><a href=\"/tasks?status=waiting\">Waiting For</a></li><li><a href=\"/projects\" class=\"text-primary font-medium\">Projects</a></li><li><a href=\"/tasks?status=someday\">Someday/Maybe</a></li><li><a href=\"/weekly-review\" class=\"text-accent\">Weekly Review</a></li></ul></div></header><main><div id=\"main-content\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div id=\"search-results\" class=\"mt-4\"></div></main><footer class=\"footer footer-center p-4 bg-base-100 text-base-content mt-6 rounded-box\"><div><p>Copyright © 2025 - All rights reserved</p></div></footer></div><!-- Quick capture modal -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.QuickCaptureModal().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "<!-- Custom JS --><script src=\"/static/js/main.js\"></script><script>\n\t\t\t// Function to submit quick capture form with Ctrl+Enter\n\t\t\tfunction submitQuickCapture() {\n\t\t\t\tdocument.getElementById('quick-capture-submit').click();\n\t\t\t}\n\t\t\t\n\t\t\t// Global keyboard shortcut for quick capture (ALT+N)\n\t\t\tdocument.addEventListener('keydown', function(e) {\n\t\t\t\tif (e.altKey && e.key === 'n') {\n\t\t\t\t\te.preventDefault();\n\t\t\t\t\tdocument.getElementById('quick-capture-modal').showModal();\n\t\t\t\t}\n\t\t\t});\n\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
