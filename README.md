# trm-cli
trm cli

Run following CLI commands:
1.	Login :  trm login ( enter name username and password )
2.	Health : trm health
3.	Set SAP server URL: trm set url <sap-url>
4.	Read SAP users: trm sap-user read all or trm sap-user read <username>
5.	Create SAP user: trm sap-user create <username>
6.	Clone SAP user : trm sap-user clone <usernameCloneFrom> <NewUsername>
7.	Unit test: trm unit-test <package>
8.	Code inspector: trm qa-check <”package” or “object”> <value>
9.	Read git repository: trm repo read
10.	Set git repository: trm repo set <repoName>
11.	Git add: trm git add <”.” Or “objectname”>
    
  11.1 Enter GitHub username: <username>
  11.2 Enter GitHub password: <password>
  
12.	Git push : trm git push
    
  12.1 Enter GitHub password: <password>
  12.2 Enter author name: 
  12.3 Enter author email:
  12.4 Enter comment:
