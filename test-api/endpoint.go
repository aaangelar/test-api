package user

import()


//declare Endpoint struct
type Endpoint struct {
	//{fucntionName}Endpoint  
	GetUsersEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetExportTypeTemplateEndpoint:         MakeGetExportTypeTemplateEndpoint(s),
	}
}

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		GetExportTypeTemplateEndpoint:         httptransport.NewClient("POST", tgt, encodeGetUsersRequest, 
			decodeGetUsersResponse, options...).Endpoint(),
	}, nil
}


//declare Make[function]Endpoint
func MakeGetUsersEnpoint (s Service) endpoint.Endpoint {
	return func (ctx context.Context, r *http.Request) (interface{}, error) {
		req := request.(GetUsersRequest)
		p, e := s.GetUsers(ctx, req.UserID)

		return GetUsersResponse{
			Details: p,
		}, nil
	}
}

func (r GetExportTypeTemplateResponse) error() error { return r.Err }

