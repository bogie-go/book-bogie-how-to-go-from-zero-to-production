// Shared admin shell and resource pages. Every admin page includes this
// file, then calls adminReady() and adminResource(). The shop's pages
// keep their own script; this one is only for /admin.
//
// The gate is a convenience, not the security: every /admin API call is
// checked again by the server (RequireStaff). Skipping the redirect
// gets a visitor nothing.

const tokenKey = "depot_token";
const userKey = "depot_user";

const getToken = () => localStorage.getItem(tokenKey);

function currentUser() {
  try { return JSON.parse(localStorage.getItem(userKey) || "null"); } catch (e) { return null; }
}

const esc = (t) => String(t).replace(/[&<>"']/g, (c) => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;"
}[c]));

async function api(method, path, body, retried) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  const route = path.split("?")[0];
  const needsToken = route !== "/products" && !route.startsWith("/products/") &&
    !route.startsWith("/categories") && !route.startsWith("/guest_sessions") &&
    !route.startsWith("/sessions") && route !== "/users";
  if (needsToken) {
    const token = getToken();
    if (!token) {
      const err = new Error("not signed in");
      err.status = 401;
      throw err;
    }
    headers["Authorization"] = "Bearer " + token;
  } else if (getToken()) {
    headers["Authorization"] = "Bearer " + getToken();
  }
  const res = await fetch(path, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  if (res.status === 401 && needsToken && !retried) {
    localStorage.removeItem(tokenKey);
    localStorage.removeItem(userKey);
    const err = new Error("not signed in");
    err.status = 401;
    throw err;
  }
  let data = null;
  try { data = await res.json(); } catch (e) {}
  if (!res.ok) {
    const err = new Error((data && data.error) || ("HTTP " + res.status));
    err.status = res.status;
    throw err;
  }
  return data;
}

function adminLinks(includeCategories) {
  const links = [
    { label: "Products", href: "/admin/products.html" },
    { label: "Sales", href: "/admin/price_rules.html" },
    { label: "Coupons", href: "/admin/coupons.html" },
  ];
  if (includeCategories) links.push({ label: "Categories", href: "/admin/categories.html" });
  return links;
}

function logout() {
  localStorage.removeItem(tokenKey);
  localStorage.removeItem(userKey);
  location.href = "/";
}

// Fill the sidebar, check that the visitor is staff, and reveal the shell.
// Returns true if the page should load its data. Not signed in goes to
// /signin.html; a customer sees "Staff only".
async function adminReady(opts) {
  opts = opts || {};
  const nav = document.getElementById("admin-nav");
  const here = location.pathname;
  nav.innerHTML = "";
  for (const link of adminLinks(!!opts.categories)) {
    const a = document.createElement("a");
    a.href = link.href;
    a.textContent = link.label;
    if (here === link.href) a.className = "active";
    nav.appendChild(a);
  }
  const back = document.createElement("a");
  back.href = "/";
  back.innerHTML = "&larr; Back to the shop";
  nav.appendChild(back);
  document.getElementById("signout").addEventListener("click", logout);

  const gate = document.getElementById("gate");
  const user = currentUser();
  if (!user || !getToken()) {
    location.href = "/signin.html";
    return false;
  }
  document.getElementById("admin-email").textContent = user.email || "";
  try {
    await api("GET", "/admin/products");
    user.staff = true;
    localStorage.setItem(userKey, JSON.stringify(user));
    document.getElementById("shell").hidden = false;
    gate.hidden = true;
    return true;
  } catch (err) {
    if (err.status === 401) {
      localStorage.removeItem(tokenKey);
      localStorage.removeItem(userKey);
      location.href = "/signin.html";
      return false;
    }
    gate.textContent = err.status === 403
      ? "Staff only. Sign in with a staff account."
      : err.message;
    return false;
  }
}

function toInput(type, v) {
  if (v == null || v === "") return "";
  if (type === "json") return JSON.stringify(v);
  if (type === "datetime") {
    const d = new Date(v);
    return isNaN(d) ? "" : d.toISOString().slice(0, 16);
  }
  return v;
}

function fromInput(type, v) {
  if (type === "number") return Number(v);
  if (type === "json") return JSON.parse(v || "{}");
  if (type === "datetime") return v ? new Date(v).toISOString() : null;
  return v;
}

function cellText(type, v) {
  if (v == null || v === "") return "";
  if (type === "json") return JSON.stringify(v);
  return v;
}

// Table plus a New/Edit form on the same page, driven by RES (slug, key,
// fields). opts.archive adds Archive/Restore on products; opts.productCategories
// adds the category checkboxes from Chapter 13.
function adminResource(RES, opts) {
  opts = opts || {};
  let editing = null;
  let catsEditable = true;

  function buildForm() {
    const form = document.getElementById("form");
    form.innerHTML = RES.fields.map(([name, label, type]) => {
      const field = type === "json"
        ? '<textarea name="' + name + '" rows="4"></textarea>'
        : '<input name="' + name + '" ' +
          (type === "number" ? 'type="number"' : type === "datetime" ? 'type="datetime-local"' : 'type="text"') +
          ">";
      return '<label class="wide">' + esc(label) + field + "</label>";
    }).join("") +
      (opts.productCategories ? '<fieldset id="cats"><legend>Categories</legend><span id="cat-boxes"></span> <span class="note" id="cat-note"></span></fieldset>' : "") +
      '<div class="actions"><button type="submit" class="primary">Save</button> <button type="button" class="secondary" id="reset">Clear</button></div>';
    form.querySelector("#reset").addEventListener("click", () => startEdit(null));
  }

  function startEdit(item) {
    editing = item;
    document.getElementById("form-title").textContent = item ? "Edit #" + item.id : "New";
    for (const [name, , type] of RES.fields) {
      const el = document.querySelector('[name="' + name + '"]');
      el.value = item ? toInput(type, item[name]) : "";
    }
    document.getElementById("message").textContent = "";
    if (opts.productCategories) loadCats(item);
  }

  async function load() {
    const status = document.getElementById("status");
    status.className = "note";
    try {
      const data = await api("GET", "/admin/" + RES.slug);
      const items = data[RES.key] || [];
      status.textContent = items.length ? "" : "Nothing here yet.";
      document.querySelector("#list thead").innerHTML = "<tr><th>ID</th>" +
        RES.fields.map((f) => "<th>" + esc(f[1].split(" (")[0]) + "</th>").join("") +
        "<th></th></tr>";
      const body = document.querySelector("#list tbody");
      body.innerHTML = "";
      for (const item of items) {
        const tr = document.createElement("tr");
        let actions = '<a href="#" class="edit">Edit</a> <a href="#" class="delete">Delete</a>';
        if (opts.archive && RES.slug === "products") actions += ' <a href="#" class="archive"></a>';
        tr.innerHTML = "<td></td>" + RES.fields.map(() => "<td></td>").join("") +
          '<td class="actions">' + actions + "</td>";
        const cells = tr.querySelectorAll("td");
        cells[0].textContent = item.id;
        RES.fields.forEach(([name, , type], i) => {
          cells[i + 1].textContent = cellText(type, item[name]);
        });
        tr.querySelector(".edit").addEventListener("click", (e) => {
          e.preventDefault();
          startEdit(item);
        });
        tr.querySelector(".delete").addEventListener("click", async (e) => {
          e.preventDefault();
          if (!confirm("Delete #" + item.id + "?")) return;
          try {
            await api("DELETE", "/admin/" + RES.slug + "/" + item.id);
            startEdit(null);
            load();
          } catch (err) {
            document.getElementById("message").textContent = err.message;
          }
        });
        if (opts.archive && RES.slug === "products") {
          const archived = !!item.archived_at;
          if (archived) tr.classList.add("archived");
          const link = tr.querySelector(".archive");
          link.textContent = archived ? "Restore" : "Archive";
          link.addEventListener("click", async (e) => {
            e.preventDefault();
            try {
              await api("POST", "/admin/products/" + item.id + (archived ? "/restore" : "/archive"));
              load();
            } catch (err) {
              document.getElementById("message").textContent = err.message;
            }
          });
        }
        body.appendChild(tr);
      }
      document.getElementById("list").hidden = items.length === 0;
    } catch (err) {
      status.textContent = err.status === 403
        ? "Staff only. Sign in with a staff account."
        : err.message;
      status.className = "err";
    }
  }

  async function buildCats() {
    try {
      const data = await api("GET", "/admin/categories");
      const boxes = document.getElementById("cat-boxes");
      if (!data.categories || !data.categories.length) {
        boxes.innerHTML = '<span class="note">No categories yet. Add some on the Categories page.</span>';
        return;
      }
      boxes.innerHTML = data.categories.map((c) =>
        '<label class="check"><input type="checkbox" name="category" value="' + c.id + '"> ' +
        esc(c.name) + "</label>").join("");
    } catch (e) {}
  }

  function checked() {
    return [...document.querySelectorAll('[name="category"]:checked')].map((b) => Number(b.value));
  }

  async function loadCats(item) {
    const boxes = document.querySelectorAll('[name="category"]');
    boxes.forEach((b) => { b.checked = false; b.disabled = false; });
    catsEditable = true;
    const note = document.getElementById("cat-note");
    if (note) note.textContent = "";
    if (!item) return;
    try {
      const shown = await api("GET", "/products/" + item.id);
      const mine = new Set((shown.categories || []).map((c) => c.id));
      boxes.forEach((b) => { b.checked = mine.has(Number(b.value)); });
    } catch (e) {
      catsEditable = false;
      boxes.forEach((b) => { b.disabled = true; });
      if (note) note.textContent = "Restore this product to change its categories.";
    }
  }

  async function saveCats(id) {
    if (catsEditable) await api("PUT", "/admin/products/" + id + "/categories", { category_ids: checked() });
  }

  document.getElementById("form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const message = document.getElementById("message");
    message.textContent = "";
    try {
      const body = {};
      for (const [name, , type] of RES.fields) {
        body[name] = fromInput(type, document.querySelector('[name="' + name + '"]').value);
      }
      const saved = editing
        ? await api("PATCH", "/admin/" + RES.slug + "/" + editing.id, body)
        : await api("POST", "/admin/" + RES.slug, body);
      if (opts.productCategories) await saveCats(saved.id);
      startEdit(null);
      load();
    } catch (err) {
      message.textContent = err.message;
    }
  });

  buildForm();
  if (opts.productCategories) buildCats();
  startEdit(null);
  load();
}
