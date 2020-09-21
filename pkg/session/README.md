# session package

This library provides a simple and consistent REST over HTTP library for accessing Akamai Endpoints

## Depedencies 

This library is dependent on the `github.com/akamai/AkamaiOPEN-edgegrid-golang/pkg/edgegrid` interface.

## Basic Example

```
func main() {
     edgerc := Must(New())

     s, err := session.New(
         session.WithConfig(edgerc),
     )
     if err != nil {
         panic(err)
     }

     var contracts struct {
		AccountID string         `json:"accountId"`
		Contracts ContractsItems `json:"contracts"`
        Items []struct {
            ContractID       string `json:"contractId"`
		    ContractTypeName string `json:"contractTypeName"`
        } `json:"items"`
     }

     req, _ := http.NewRequest(http.MethodGet, "/papi/v1/contracts", nil)

     _, err := s.Exec(r, &contracts)
     if err != nil {
         panic(err);
     }

     // do something with contracts
}
        
```