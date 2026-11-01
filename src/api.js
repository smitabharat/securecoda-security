const API = process.env.REACT_APP_API || 'http://localhost:8080/api'
export async function fetchAlerts(){
const res = await fetch(API + '/alerts')
return res.json()
}
export async function unpublish(docId){
return fetch(API + '/remediate/unpublish/' + docId, {method:'POST'})
}
export async function deleteDoc(docId){
return fetch(API + '/remediate/delete-doc/' + docId, {method:'POST'})
}