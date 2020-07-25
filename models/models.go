package models

type BackendUrl struct {
	// This path will be concatenated with the `PortalApplication`'s name. e.g. 
	// if PortalApplication name is `bcat` and path is /transactions the resultant
	// web server config that will be generated will be:
	// 
	// location /bcat/transactions {
  	//	proxy_pass <endpoint>/bcat/transactions
	// }
	path string

	// The scheme://<hostname>:<port>
	backendHost string

	sslVerify bool

	// TODO: SSL Verification

}

type PortalApplication struct {
	// the name of the portal application, usually equivalent to the app code
	Name string

	// a user-friendly way of referring to the portal application
	Label string

	// the name of the npm package where the application's user interface is
	// UIPackage string

	// Various endpoints for various paths
	// BackendUrls []BackendUrl
}
