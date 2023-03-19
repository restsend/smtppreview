const serveraddr = 'http://localhost:8000'
export async function getSummary() {
    let req = await fetch(`${serveraddr}/api/summary`, { method: 'POST' })
    return await req.json()
}

export async function listMails(pos, keyword, limit = 20) {
    let req = await fetch(`${serveraddr}/api/mail/query`, {
        method: 'POST',
        body: JSON.stringify({
            pos, keyword, limit
        })
    })
    return await req.json()
}

export async function markMailOpened(msgid, opened) {
    let req = await fetch(`${serveraddr}/api/mail/${msgid}`, {
        method: 'PATCH',
        body: JSON.stringify({
            opened
        })
    })
    return await req.json()
}

export async function deleteMail(msgid) {
    let req = await fetch(`${serveraddr}/api/mail/${msgid}`, {
        method: 'DELETE'
    })
    return await req.json()
}