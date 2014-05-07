[1mdiff --git a/mailpin.go b/mailpin.go[m
[1mindex 2d16a96..e98c930 100644[m
[1m--- a/mailpin.go[m
[1m+++ b/mailpin.go[m
[36m@@ -33,7 +33,8 @@[m [mfunc (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {[m
 [m
 func init() {[m
 	r := mux.NewRouter()[m
[31m-	r.Handle("/", appHandler(root))[m
[32m+[m	[32mr.Handle("/", appHandler(root)).Methods("GET")[m
[32m+[m	[32mr.Handle("/{id}", appHandler(viewEmail)).Methods("GET")[m
 	http.Handle("/", r)[m
 	http.Handle("/_ah/mail/", appHandler(incomingMail))[m
 }[m
[36m@@ -46,6 +47,20 @@[m [mMail to p@go-mailpin.appspotmail.com. Get a short sharable URL.[m
 	return nil[m
 }[m
 [m
[32m+[m[32mfunc viewEmail(w http.ResponseWriter, r *http.Request) *appError {[m
[32m+[m	[32mc := appengine.NewContext(r)[m
[32m+[m	[32mvars := mux.Vars(r)[m
[32m+[m	[32mid := vars["id"][m
[32m+[m
[32m+[m	[32mpage, err := model.GetPage(c, id)[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mreturn &appError{err, "Page not found", http.StatusNotFound}[m
[32m+[m	[32m}[m
[32m+[m
[32m+[m	[32mw.Write(page.Body)[m
[32m+[m	[32mreturn nil[m
[32m+[m[32m}[m
[32m+[m
 func incomingMail(w http.ResponseWriter, r *http.Request) *appError {[m
 	c := appengine.NewContext(r)[m
 	defer r.Body.Close()[m
