async function j(url, opts) {
  const r = await fetch(url, opts);
  const t = await r.text();
  let data;
  try { data = JSON.parse(t); } catch { data = t; }
  if (!r.ok) throw new Error((data && data.error) || t || r.statusText);
  return data;
}

async function refresh() {
  const st = await j('/api/stats');
  document.getElementById('stats-line').textContent =
    `条目 ${st.Entries} · 吊销 ${st.Revoked} · 导入 ${st.Imports} · 扫描 ${st.Scans}`;
  const data = await j('/api/certs');
  const ul = document.getElementById('cert-list');
  ul.innerHTML = '';
  (data.certs || []).forEach(c => {
    const li = document.createElement('li');
    li.innerHTML = `<strong>${c.name}</strong> <span class="muted">${c.id}</span><br/>
      <span class="${c.revoked ? 'bad' : 'ok'}">${c.revoked ? '已吊销' : (c.active ? '活动' : '非活动')}</span>
      · 到期 ${c.not_after}`;
    ul.appendChild(li);
  });
}

document.getElementById('btn-refresh').onclick = () => refresh().catch(e => alert(e.message));

document.getElementById('import-form').onsubmit = async (ev) => {
  ev.preventDefault();
  const out = document.getElementById('import-out');
  try {
    const body = {
      name: document.getElementById('imp-name').value,
      cert_pem: document.getElementById('imp-cert').value,
      key_pem: document.getElementById('imp-key').value,
      active: true,
    };
    const res = await j('/api/certs/import', {
      method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body)
    });
    out.textContent = '导入成功 id=' + res.id;
    await refresh();
  } catch (e) { out.textContent = e.message; out.className = 'out bad'; }
};

document.getElementById('scan-form').onsubmit = async (ev) => {
  ev.preventDefault();
  const out = document.getElementById('scan-out');
  try {
    const days = Number(document.getElementById('scan-days').value || 30);
    const res = await j('/api/scan', {
      method: 'POST', headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ within_days: days })
    });
    out.textContent = JSON.stringify(res.hits || [], null, 2);
  } catch (e) { out.textContent = e.message; }
};

document.getElementById('preview-form').onsubmit = async (ev) => {
  ev.preventDefault();
  const out = document.getElementById('preview-out');
  try {
    const body = {
      old_id: document.getElementById('prev-id').value,
      cert_pem: document.getElementById('prev-cert').value,
      key_pem: document.getElementById('prev-key').value,
    };
    const res = await j('/api/rotate/preview', {
      method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body)
    });
    out.textContent = JSON.stringify(res, null, 2);
  } catch (e) { out.textContent = e.message; }
};

refresh().catch(() => {});
