package github_standards

const templateName string = "github-standards"

// func decorators(re *github_standards.GHSResponse, conf *config.Config, navItem *navigation.NavigationItem, r *http.Request) {
// 	re.Organisation = conf.Organisation
// 	re.PageTitle = navItem.Name
// 	if len(conf.Navigation) > 0 {
// 		top, active := navigation.Level(conf.Navigation, r)
// 		re.NavigationActive = active
// 		re.NavigationTop = top
// 		re.NavigationSide = active.Navigation
// 	}
// }

// func ListHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
// 	var data interface{}
// 	mapData := map[string]interface{}{}
// 	if responses, err := getter.ApiHttpResponses(navItem, r); err == nil {
// 		count := len(responses)
// 		for key, rep := range responses {
// 			gh, err := convert.UnmarshalR[*github_standards.GHSResponse](rep)
// 			if err != nil {
// 				return
// 			}
// 			// set the nav and org details
// 			decorators(gh, conf, navItem, r)
// 			if count > 1 {
// 				mapData[key] = gh
// 				data = mapData
// 			} else {
// 				data = gh
// 			}
// 		}
// 	}
// 	helpers.OutputHandler(templates, navItem.Template, data, w)

// }

// func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
// 	nav := conf.Navigation
// 	navItems := navigation.ForTemplateList(templateName, nav)
// 	for _, navItem := range navItems {
// 		var handler = func(w http.ResponseWriter, r *http.Request) {
// 			ListHandler(w, r, templates, conf, navItem)
// 		}
// 		slog.Info("[front] register", slog.String("endpoint", "githug_standards"), slog.String("list", navItem.Uri))
// 		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
// 	}
// }
