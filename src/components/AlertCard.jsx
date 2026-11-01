import React from 'react'
export default function AlertCard({alert, onUnpublish, onDelete}){
return (
<div className="border p-3 rounded">
<div className="flex justify-between">
<div>
<div className="font-bold">{alert.title || alert.doc_id}</div>
<div className="text-sm">Type: {alert.type} — {alert.detail}</div>
</div>
<div className="space-x-2">
<button onClick={onUnpublish} className="px-2 py-1 border rounded">Unpublish</button>
<button onClick={onDelete} className="px-2 py-1 border rounded">Delete</button>
</div>
</div>
</div>
)
}