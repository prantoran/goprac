* Viewing Static Single Assignment
GOSSAFUNC=main GOOS=linux GOARCH=amd64 go build -gcflags “-S” simple.go


When you open ssa.html, a number of passes will be shown, most of which are collapsed. The start pass is the SSA that is generated from the AST; the lower pass converts the non-machine specific SSA to machine-specific SSA and genssa is the final generated machine code.
