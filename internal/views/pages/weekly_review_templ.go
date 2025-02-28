// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/melihkorkmaz/gtd/internal/views/layouts"

func WeeklyReviewPage() templ.Component {
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
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"grid gap-6\"><div class=\"card bg-base-100 shadow-lg\"><div class=\"card-body\"><h2 class=\"card-title text-2xl mb-4\">Weekly Review</h2><p class=\"mb-4\">The weekly review is a time to get clear, get current, and get creative. Use this checklist to guide your weekly review process.</p><div class=\"flex justify-end mb-4\"><button class=\"btn btn-primary\" id=\"start-review\">Start Weekly Review</button> <button class=\"btn btn-success ml-2\" id=\"reset-review\" style=\"display: none;\">Reset Review</button></div><div id=\"review-progress\" class=\"mb-4\" style=\"display: none;\"><progress class=\"progress progress-primary w-full\" id=\"review-progress-bar\" value=\"0\" max=\"100\"></progress><p class=\"text-center mt-2\"><span id=\"review-progress-text\">0%</span> complete</p></div><div id=\"review-steps\" style=\"display: none;\"><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>1. Collect Loose Papers and Materials</span></div><div class=\"collapse-content\"><p>Gather all physical materials - notes, receipts, documents, business cards, etc. - into your inbox for processing.</p><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Collected all physical materials</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>2. Process Your Notes</span></div><div class=\"collapse-content\"><p>Go through any paper or digital notes you've taken during the week and transfer them to the appropriate system.</p><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Processed all notes</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>3. Empty Your Inbox</span></div><div class=\"collapse-content\"><p>Process all items in your inbox to zero. Decide what each item is and what needs to be done with it.</p><a href=\"/tasks?status=inbox\" target=\"_blank\" class=\"btn btn-outline btn-sm mt-2\">Go to Inbox</a><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Inbox is empty</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>4. Review Next Actions Lists</span></div><div class=\"collapse-content\"><p>Review your Next Actions list. Mark completed items as done and update any that have changed.</p><a href=\"/tasks?status=next\" target=\"_blank\" class=\"btn btn-outline btn-sm mt-2\">Review Next Actions</a><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Next Actions list is current and complete</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>5. Review Waiting For List</span></div><div class=\"collapse-content\"><p>Review items you're waiting on from others. Record any necessary follow-ups.</p><a href=\"/tasks?status=waiting\" target=\"_blank\" class=\"btn btn-outline btn-sm mt-2\">Review Waiting Items</a><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Waiting For list is up to date</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>6. Review Projects</span></div><div class=\"collapse-content\"><p>Review the status of all current projects. Ensure each has at least one next action.</p><a href=\"/projects\" target=\"_blank\" class=\"btn btn-outline btn-sm mt-2\">Review Projects</a><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">All projects have clear next actions</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>7. Review Someday/Maybe List</span></div><div class=\"collapse-content\"><p>Review your Someday/Maybe items. Move any to active projects if you're ready to start them.</p><a href=\"/tasks?status=someday\" target=\"_blank\" class=\"btn btn-outline btn-sm mt-2\">Review Someday/Maybe</a><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Someday/Maybe list is reviewed</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div><div x-data=\"{ open: false, completed: false }\" class=\"collapse collapse-arrow bg-base-200 mb-4\"><input type=\"checkbox\" x-bind:checked=\"open\" @click=\"open = !open\"><div class=\"collapse-title text-xl font-medium flex items-center\"><input type=\"checkbox\" class=\"checkbox mr-3\" x-model=\"completed\" @change=\"updateProgress()\"> <span>8. Get Creative</span></div><div class=\"collapse-content\"><p>Consider new ideas, possibilities, or projects you might want to pursue.</p><div class=\"form-control mt-2\"><label class=\"cursor-pointer label\"><span class=\"label-text\">Considered new ideas and possibilities</span> <input type=\"checkbox\" class=\"checkbox checkbox-primary\"></label></div></div></div></div><div id=\"review-complete\" class=\"alert alert-success shadow-lg\" style=\"display: none;\"><div><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"stroke-current flex-shrink-0 h-6 w-6\" fill=\"none\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z\"></path></svg><div><span class=\"font-bold\">Weekly review completed!</span><p class=\"text-sm\">Your GTD system is now current and up to date.</p></div></div></div></div></div></div><script>\n\t\t\t// Weekly review functionality\n\t\t\tdocument.addEventListener('DOMContentLoaded', function() {\n\t\t\t\tconst startButton = document.getElementById('start-review');\n\t\t\t\tconst resetButton = document.getElementById('reset-review');\n\t\t\t\tconst reviewSteps = document.getElementById('review-steps');\n\t\t\t\tconst reviewProgress = document.getElementById('review-progress');\n\t\t\t\tconst reviewComplete = document.getElementById('review-complete');\n\t\t\t\tconst progressBar = document.getElementById('review-progress-bar');\n\t\t\t\tconst progressText = document.getElementById('review-progress-text');\n\t\t\t\t\n\t\t\t\tstartButton.addEventListener('click', function() {\n\t\t\t\t\treviewSteps.style.display = 'block';\n\t\t\t\t\treviewProgress.style.display = 'block';\n\t\t\t\t\tstartButton.style.display = 'none';\n\t\t\t\t\tresetButton.style.display = 'inline-flex';\n\t\t\t\t});\n\t\t\t\t\n\t\t\t\tresetButton.addEventListener('click', function() {\n\t\t\t\t\tlocation.reload();\n\t\t\t\t});\n\t\t\t\t\n\t\t\t\t// Function to update progress\n\t\t\t\twindow.updateProgress = function() {\n\t\t\t\t\tconst steps = document.querySelectorAll('#review-steps > div');\n\t\t\t\t\tlet completed = 0;\n\t\t\t\t\t\n\t\t\t\t\tsteps.forEach(step => {\n\t\t\t\t\t\tif (step.__x && step.__x.$data.completed) {\n\t\t\t\t\t\t\tcompleted++;\n\t\t\t\t\t\t}\n\t\t\t\t\t});\n\t\t\t\t\t\n\t\t\t\t\tconst percentage = Math.round((completed / steps.length) * 100);\n\t\t\t\t\tprogressBar.value = percentage;\n\t\t\t\t\tprogressText.textContent = percentage + '%';\n\t\t\t\t\t\n\t\t\t\t\tif (percentage === 100) {\n\t\t\t\t\t\treviewComplete.style.display = 'block';\n\t\t\t\t\t} else {\n\t\t\t\t\t\treviewComplete.style.display = 'none';\n\t\t\t\t\t}\n\t\t\t\t};\n\t\t\t});\n\t\t</script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = layouts.Base("Weekly Review - GTD App").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
