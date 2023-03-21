
export async function getSummary() {
    let req = await fetch(`/api/summary`, { method: 'POST' })
    return await req.json()
}

export async function listMails(pos, keyword, limit = 20) {
    let req = await fetch(`/api/mail/query`, {
        method: 'POST',
        body: JSON.stringify({
            pos, keyword, limit,
            orders: [{ name: 'createdAt', op: 'desc' }]
        })
    })
    return await req.json()
}

export async function markMailOpened(msgid, opened) {
    let req = await fetch(`/api/mail/${msgid}`, {
        method: 'PATCH',
        body: JSON.stringify({
            opened
        })
    })
    return await req.json()
}

export async function deleteMail(msgid) {
    let req = await fetch(`/api/mail/${msgid}`, {
        method: 'DELETE'
    })
    return await req.json()
}



export function hasAttachment(msg) {
    if (msg.attachments) {
        let files = JSON.parse(msg.attachments) || []
        return files.length > 0
    }
}