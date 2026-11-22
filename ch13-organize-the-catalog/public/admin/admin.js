// Shared admin shell, list pages, and form pages. Every admin page
// includes this file, then calls adminReady() and one of adminDashboard,
// adminList, or adminForm. The shop's pages keep their own script.
//
// The gate is a convenience, not the security: every /admin API call is
// checked again by the server (RequireStaff). Skipping the redirect
// gets a visitor nothing.
//
// Sidebar labels and "New …" button text live in ADMIN_SECTIONS.
// Promotions pages are promotions.html / promotion.html; the API
// they call stays /admin/price_rules.

const tokenKey = "depot_token";
const userKey = "depot_user";

const ADMIN_SECTIONS = {
  dashboard: {
    label: "Dashboard",
    href: "/admin/index.html",
  },
  products: {
    label: "Products",
    list: "/admin/products.html",
    form: "/admin/product.html",
    newLabel: "New product",
    slug: "products",
    key: "products",
    columns: [
      ["id", "ID"],
      ["name", "Name"],
      ["price_cents", "Price"],
      ["stock", "Stock"],
    ],
    fields: [
      ["name", "Name", "text", "wide"],
      ["description", "Description", "textarea", "wide"],
      ["price_cents", "Price (cents)", "number"],
      ["stock", "Stock", "number"],
      ["image_url", "Image URL", "text", "wide"],
    ],
  },
  // User-facing name and files are Promotions; the API stays price_rules.
  promotions: {
    label: "Promotions",
    list: "/admin/promotions.html",
    form: "/admin/promotion.html",
    newLabel: "New promotion",
    slug: "price_rules",
    key: "price_rules",
    columns: [
      ["id", "ID"],
      ["name", "Name"],
      ["kind", "Kind"],
      ["priority", "Priority"],
      ["starts_at", "Starts at"],
      ["ends_at", "Ends at"],
    ],
    fields: [
      ["name", "Name", "text", "wide"],
      ["kind", "Kind (percentage_discount, buy_x_get_y, minimum_purchase)", "text", "wide"],
      ["properties", "Properties (JSON)", "json", "wide"],
      ["priority", "Priority", "number"],
      ["starts_at", "Starts at", "datetime"],
      ["ends_at", "Ends at", "datetime"],
    ],
  },
  coupons: {
    label: "Coupons",
    list: "/admin/coupons.html",
    form: "/admin/coupon.html",
    newLabel: "New coupon",
    slug: "coupons",
    key: "coupons",
    columns: [
      ["id", "ID"],
      ["code", "Code"],
      ["price_rule_id", "Price rule id"],
    ],
    fields: [
      ["code", "Code (upper case)", "text", "wide"],
      ["price_rule_id", "Price rule id", "number"],
    ],
  },
  categories: {
    label: "Categories",
    list: "/admin/categories.html",
    form: "/admin/category.html",
    newLabel: "New category",
    slug: "categories",
    key: "categories",
    columns: [
      ["id", "ID"],
      ["name", "Name"],
    ],
    fields: [
      ["name", "Name", "text", "wide"],
    ],
  },
};

const getToken = () => localStorage.getItem(tokenKey);

function currentUser() {
  try { return JSON.parse(localStorage.getItem(userKey) || "null"); } catch (e) { return null; }
}

const esc = (t) => String(t).replace(/[&<>"']/g, (c) => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;"
}[c]));

function queryId() {
  return new URLSearchParams(location.search).get("id");
}

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

function adminNavItems(includeCategories) {
  const items = [
    ADMIN_SECTIONS.products,
    ADMIN_SECTIONS.promotions,
    ADMIN_SECTIONS.coupons,
  ];
  if (includeCategories) items.push(ADMIN_SECTIONS.categories);
  return items;
}

function isActiveSection(section) {
  const here = location.pathname;
  if (section.href && here === section.href) return true;
  if (section.list && here === section.list) return true;
  if (section.form && here === section.form) return true;
  return false;
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
  nav.innerHTML = "";
  for (const section of adminNavItems(!!opts.categories)) {
    const a = document.createElement("a");
    a.href = section.href || section.list;
    a.textContent = section.label;
    if (isActiveSection(section)) a.className = "active";
    nav.appendChild(a);
  }
  const back = document.createElement("a");
  back.href = "/";
  back.innerHTML = "&larr; Back to the shop";
  nav.appendChild(back);
  const brand = document.querySelector(".sidebar .brand");
  if (brand) brand.href = ADMIN_SECTIONS.dashboard.href;
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
  if (type === "json") return JSON.stringify(v, null, 2);
  if (type === "datetime") {
    const d = new Date(v);
    if (isNaN(d) || d.getUTCFullYear() < 1000) return "";
    const p = (n) => String(n).padStart(2, "0");
    return d.getFullYear() + "-" + p(d.getMonth() + 1) + "-" + p(d.getDate()) +
      "T" + p(d.getHours()) + ":" + p(d.getMinutes());
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
  if (type === "datetime" || (typeof v === "string" && /^\d{4}-\d{2}-\d{2}T/.test(v))) {
    const d = new Date(v);
    if (isNaN(d) || d.getUTCFullYear() < 1000) return "";
    return d.toLocaleString();
  }
  return v;
}

function fieldType(section, name) {
  const field = (section.fields || []).find((f) => f[0] === name);
  return field ? field[2] : "text";
}

function itemLabel(item) {
  return item.name || item.code || ("#" + item.id);
}

// Dashboard cards, one per section, with a live count from the list API.
async function adminDashboard(opts) {
  opts = opts || {};
  const cards = document.getElementById("cards");
  const status = document.getElementById("status");
  const sections = adminNavItems(!!opts.categories).filter((s) => s.slug);
  cards.innerHTML = "";
  for (const section of sections) {
    const card = document.createElement("div");
    card.className = "card";
    card.innerHTML = "<h2><a href=\"" + esc(section.list) + "\">" + esc(section.label) +
      "</a></h2><p class=\"count\">Loading&hellip;</p>";
    cards.appendChild(card);
    try {
      const data = await api("GET", "/admin/" + section.slug);
      const n = (data[section.key] || []).length;
      card.querySelector(".count").textContent = n + (n === 1 ? " record" : " records");
    } catch (err) {
      card.querySelector(".count").textContent = err.message;
    }
  }
  if (status) status.textContent = "";
}

// List page: a "New …" button, a table, themed Edit/Delete (and
// Archive/Restore on products from Chapter 12).
function adminList(section, opts) {
  opts = opts || {};
  document.title = section.label + " - Bogie Depot Admin";
  document.getElementById("list-title").textContent = section.label;
  const newBtn = document.getElementById("new-link");
  newBtn.href = section.form;
  newBtn.textContent = section.newLabel;

  const status = document.getElementById("status");
  const message = document.getElementById("message");
  const table = document.getElementById("list");

  async function load() {
    status.className = "note";
    message.textContent = "";
    try {
      const data = await api("GET", "/admin/" + section.slug);
      const items = data[section.key] || [];
      status.textContent = items.length ? "" : "Nothing here yet.";
      table.querySelector("thead").innerHTML = "<tr>" +
        section.columns.map((c) => "<th>" + esc(c[1]) + "</th>").join("") +
        "<th></th></tr>";
      const body = table.querySelector("tbody");
      body.innerHTML = "";
      for (const item of items) {
        const tr = document.createElement("tr");
        let actions = '<a class="btn secondary" href="' + section.form + "?id=" +
          encodeURIComponent(item.id) + '">Edit</a> ' +
          '<button type="button" class="danger delete">Delete</button>';
        if (opts.archive) {
          actions += ' <button type="button" class="secondary archive"></button>';
        }
        tr.innerHTML = section.columns.map(() => "<td></td>").join("") +
          '<td class="actions">' + actions + "</td>";
        const cells = tr.querySelectorAll("td");
        section.columns.forEach(([name], i) => {
          if (name === "name" && opts.archive && item.archived_at) {
            cells[i].innerHTML = esc(item[name] || "") +
              ' <span class="badge">Archived</span>';
          } else {
            cells[i].textContent = name === "id" ? item.id :
              cellText(fieldType(section, name), item[name]);
          }
        });
        if (opts.archive && item.archived_at) tr.classList.add("archived");
        tr.querySelector(".delete").addEventListener("click", async () => {
          if (!confirm("Delete \"" + itemLabel(item) + "\"?")) return;
          try {
            await api("DELETE", "/admin/" + section.slug + "/" + item.id);
            load();
          } catch (err) {
            message.textContent = err.message;
          }
        });
        if (opts.archive) {
          const archived = !!item.archived_at;
          const btn = tr.querySelector(".archive");
          btn.textContent = archived ? "Restore" : "Archive";
          btn.addEventListener("click", async () => {
            try {
              await api("POST", "/admin/products/" + item.id +
                (archived ? "/restore" : "/archive"));
              load();
            } catch (err) {
              message.textContent = err.message;
            }
          });
        }
        body.appendChild(tr);
      }
      table.hidden = items.length === 0;
    } catch (err) {
      status.textContent = err.status === 403
        ? "Staff only. Sign in with a staff account."
        : err.message;
      status.className = "err";
    }
  }

  load();
}

// Form page: /admin/product.html is a new record; ?id=5 edits product 5.
function adminForm(section, opts) {
  opts = opts || {};
  const id = queryId();
  // "New product" → "Edit product".
  const editTitle = section.newLabel.replace(/^New /, "Edit ");
  document.title = (id ? editTitle : section.newLabel) + " - Bogie Depot Admin";
  document.getElementById("form-title").textContent = id ? editTitle : section.newLabel;
  const back = document.getElementById("back-link");
  back.href = section.list;
  back.innerHTML = "&larr; " + esc(section.label);

  const form = document.getElementById("form");
  const message = document.getElementById("message");
  const missing = document.getElementById("missing");
  let catsEditable = true;

  form.innerHTML = '<div class="grid">' + section.fields.map((f) => {
    const [name, label, type, wide] = f;
    const field = type === "json" || type === "textarea"
      ? '<textarea id="f-' + name + '" name="' + name + '" rows="4"></textarea>'
      : '<input id="f-' + name + '" name="' + name + '" ' +
        (type === "number" ? 'type="number"' :
          type === "datetime" ? 'type="datetime-local"' : 'type="text"') + ">";
    return '<label class="' + (wide ? "wide" : "") + '">' + esc(label) + field + "</label>";
  }).join("") +
    (opts.productCategories
      ? '<fieldset id="cats"><legend>Categories</legend><span id="cat-boxes"></span> <span class="note" id="cat-note"></span></fieldset>'
      : "") +
    "</div>" +
    '<div class="actions"><button type="submit" class="primary">Save</button> ' +
    '<a class="btn secondary" id="cancel-in-form" href="' + section.list + '">Cancel</a></div>';

  function fill(item) {
    for (const [name, , type] of section.fields) {
      form.querySelector('[name="' + name + '"]').value =
        item ? toInput(type, item[name]) : "";
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
        '<label class="check"><input type="checkbox" name="category" value="' +
        c.id + '"> ' + esc(c.name) + "</label>").join("");
    } catch (e) {}
  }

  function checked() {
    return [...form.querySelectorAll('[name="category"]:checked')].map((b) => Number(b.value));
  }

  async function loadCats(itemId) {
    const boxes = form.querySelectorAll('[name="category"]');
    boxes.forEach((b) => { b.checked = false; b.disabled = false; });
    catsEditable = true;
    const note = document.getElementById("cat-note");
    if (note) note.textContent = "";
    if (!itemId) return;
    try {
      const shown = await api("GET", "/products/" + itemId);
      const mine = new Set((shown.categories || []).map((c) => c.id));
      boxes.forEach((b) => { b.checked = mine.has(Number(b.value)); });
    } catch (e) {
      catsEditable = false;
      boxes.forEach((b) => { b.disabled = true; });
      if (note) note.textContent = "Restore this product to change its categories.";
    }
  }

  async function saveCats(savedId) {
    if (catsEditable) {
      await api("PUT", "/admin/products/" + savedId + "/categories", {
        category_ids: checked(),
      });
    }
  }

  (async function init() {
    if (opts.productCategories) await buildCats();
    if (!id) {
      fill(null);
      return;
    }
    try {
      const item = await api("GET", "/admin/" + section.slug + "/" + encodeURIComponent(id));
      fill(item);
      if (opts.productCategories) await loadCats(id);
    } catch (err) {
      if (err.status === 404) {
        form.hidden = true;
        missing.hidden = false;
        missing.innerHTML = "No such record. <a href=\"" + section.list +
          "\">Back to " + esc(section.label) + "</a>";
        return;
      }
      message.textContent = err.message;
    }
  })();

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    message.textContent = "";
    try {
      const body = {};
      for (const [name, , type] of section.fields) {
        body[name] = fromInput(type, form.querySelector('[name="' + name + '"]').value);
      }
      const saved = id
        ? await api("PATCH", "/admin/" + section.slug + "/" + id, body)
        : await api("POST", "/admin/" + section.slug, body);
      if (opts.productCategories) await saveCats(saved.id);
      location.href = section.list;
    } catch (err) {
      message.textContent = err.message;
    }
  });
}
