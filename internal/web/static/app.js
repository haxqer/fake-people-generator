const FIELD_GROUPS = [
  {
    title: "Education",
    description: "学校、专业和毕业年份等资料。",
    items: [
      ["education.level", "Level", "学历层级"],
      ["education.degree_type", "Degree Type", "学位类型"],
      ["education.school_name", "School Name", "学校名称"],
      ["education.major", "Major", "专业"],
      ["education.graduation_year", "Graduation Year", "毕业年份"],
      ["education.student_status", "Student Status", "在读或校友状态"],
      ["education.school_city", "School City", "学校城市"],
      ["education.school_state", "School State", "学校州"],
      ["education.school_email", "School Email", "学校邮箱"],
    ],
  },
  {
    title: "Employment",
    description: "职业增强字段，适合做业务系统联调。",
    items: [
      ["employment.industry", "Industry", "所属行业"],
      ["employment.seniority", "Seniority", "职级"],
      ["employment.start_date", "Start Date", "入职日期"],
      ["employment.work_mode", "Work Mode", "办公模式"],
      ["employment.manager_name", "Manager Name", "直属主管"],
      ["employment.office_city", "Office City", "办公城市"],
      ["employment.office_state", "Office State", "办公州"],
    ],
  },
  {
    title: "Sensitive",
    description: "仅用于测试环境的合成敏感字段。",
    items: [
      ["sensitive.ssn", "SSN", "社会安全号"],
      ["sensitive.driver_license.number", "Driver License No.", "驾照号码"],
      ["sensitive.driver_license.state_code", "DL State", "驾照州代码"],
      ["sensitive.passport.number", "Passport No.", "护照号码"],
      ["sensitive.passport.country", "Passport Country", "护照国家"],
      ["sensitive.bank_card.brand", "Card Brand", "银行卡品牌"],
      ["sensitive.bank_card.number", "Card Number", "银行卡号"],
      ["sensitive.bank_card.expiry_month", "Expiry Month", "到期月份"],
      ["sensitive.bank_card.expiry_year", "Expiry Year", "到期年份"],
      ["sensitive.bank_card.cvv", "CVV", "安全码"],
    ],
  },
];

const DEFAULT_FIELDS = new Set([
  "education.level",
  "education.major",
  "education.graduation_year",
  "employment.industry",
  "employment.seniority",
  "employment.work_mode",
]);

const state = {
  activeController: null,
  autoTimer: 0,
  lastResponse: null,
  lastRequest: null,
};

const elements = {};

document.addEventListener("DOMContentLoaded", () => {
  cacheElements();
  renderFieldGroups();
  bindEvents();
  toggleAgeInputs();
  syncUI();
  generateProfiles("initial");
});

function cacheElements() {
  elements.form = document.getElementById("generator-form");
  elements.count = document.getElementById("count");
  elements.gender = document.getElementById("gender");
  elements.stateInput = document.getElementById("state");
  elements.advancedPanel = document.getElementById("advanced-panel");
  elements.advancedSummaryText = document.getElementById("advanced-summary-text");
  elements.ageEnabled = document.getElementById("age-enabled");
  elements.ageMin = document.getElementById("age-min");
  elements.ageMax = document.getElementById("age-max");
  elements.autoGenerate = document.getElementById("auto-generate");
  elements.generateButton = document.getElementById("generate-button");
  elements.resetButton = document.getElementById("reset-button");
  elements.copyRequestButton = document.getElementById("copy-request-button");
  elements.copyResponseButton = document.getElementById("copy-response-button");
  elements.downloadResponseButton = document.getElementById("download-response-button");
  elements.formMessage = document.getElementById("form-message");
  elements.fieldGroups = document.getElementById("field-groups");
  elements.fieldCount = document.getElementById("field-count");
  elements.requestPreview = document.getElementById("request-preview");
  elements.responsePreview = document.getElementById("response-preview");
  elements.metaGrid = document.getElementById("meta-grid");
  elements.profiles = document.getElementById("profiles");
  elements.statusBanner = document.getElementById("status-banner");
}

function renderFieldGroups() {
  elements.fieldGroups.innerHTML = FIELD_GROUPS.map((group) => {
    const items = group.items
      .map(([value, label, hint]) => {
        const checked = DEFAULT_FIELDS.has(value) ? "checked" : "";
        return `
          <label class="checkbox">
            <input type="checkbox" name="include_fields" value="${escapeHtml(value)}" ${checked} />
            <span>
              <strong>${escapeHtml(label)}</strong>
              <span>${escapeHtml(hint)}</span>
            </span>
          </label>
        `;
      })
      .join("");

    return `
      <section class="field-group">
        <h4>${escapeHtml(group.title)}</h4>
        <p>${escapeHtml(group.description)}</p>
        <div class="checkbox-list">${items}</div>
      </section>
    `;
  }).join("");
}

function bindEvents() {
  elements.form.addEventListener("submit", async (event) => {
    event.preventDefault();
    await generateProfiles("manual");
  });

  elements.form.addEventListener("input", () => {
    syncUI();
    scheduleAutoGenerate();
  });

  elements.form.addEventListener("change", () => {
    syncUI();
    scheduleAutoGenerate();
  });

  elements.resetButton.addEventListener("click", () => {
    resetForm();
    syncUI();
    generateProfiles("reset");
  });

  elements.copyRequestButton.addEventListener("click", async () => {
    const validationError = validateForm();
    if (validationError) {
      setFormMessage(validationError, "error");
      return;
    }

    await copyText(JSON.stringify(buildRequest(), null, 2), "请求 JSON 已复制。");
  });

  elements.copyResponseButton.addEventListener("click", async () => {
    if (!state.lastResponse) {
      setStatus("还没有可复制的响应。", "idle");
      return;
    }

    await copyText(JSON.stringify(state.lastResponse, null, 2), "响应 JSON 已复制。");
  });

  elements.downloadResponseButton.addEventListener("click", () => {
    if (!state.lastResponse) {
      setStatus("还没有可下载的响应。", "idle");
      return;
    }

    downloadJSON(state.lastResponse);
  });
}

function syncUI() {
  toggleAgeInputs();
  updateAdvancedSummary();
  updateRequestPreview();
}

function toggleAgeInputs() {
  const enabled = elements.ageEnabled.checked;
  elements.ageMin.disabled = !enabled;
  elements.ageMax.disabled = !enabled;
}

function updateAdvancedSummary() {
  const total = getSelectedFields().length;
  const ageCopy = elements.ageEnabled.checked
    ? `Age ${elements.ageMin.value}-${elements.ageMax.value}`
    : "Age off";
  const modeCopy = elements.autoGenerate.checked ? "Auto on" : "Manual";

  elements.fieldCount.textContent = `${total} ${total === 1 ? "field" : "fields"}`;
  elements.advancedSummaryText.textContent = `${ageCopy} / ${modeCopy}`;
}

function updateRequestPreview() {
  const validationError = validateForm();
  if (validationError) {
    elements.requestPreview.textContent = validationError;
    elements.generateButton.disabled = true;
    setFormMessage(validationError, "error");
    if (elements.ageEnabled.checked && validationError.includes("年龄")) {
      elements.advancedPanel.open = true;
    }
    return;
  }

  const request = buildRequest();
  state.lastRequest = request;
  elements.requestPreview.textContent = JSON.stringify(request, null, 2);
  elements.generateButton.disabled = false;
  setFormMessage("参数有效。可以手动生成，或保持自动生成开启。", "success");
}

function validateForm() {
  const count = Number(elements.count.value);
  if (!Number.isInteger(count) || count < 1 || count > 100) {
    return "Count 需要是 1 到 100 之间的整数。";
  }

  const rawState = elements.stateInput.value.trim();
  if (rawState && !/^[a-z]{2}$/i.test(rawState)) {
    return "State 需要填写 2 位州代码，例如 CA。";
  }

  if (elements.ageEnabled.checked) {
    const ageMin = Number(elements.ageMin.value);
    const ageMax = Number(elements.ageMax.value);

    if (!Number.isInteger(ageMin) || !Number.isInteger(ageMax)) {
      return "启用年龄范围后，需要同时填写最小值和最大值。";
    }
    if (ageMin < 0 || ageMin > 100 || ageMax < 0 || ageMax > 100) {
      return "年龄范围需要在 0 到 100 之间。";
    }
    if (ageMin > ageMax) {
      return "最小年龄不能大于最大年龄。";
    }
  }

  return "";
}

function buildRequest() {
  const request = {
    country: "US",
    count: Number(elements.count.value),
  };

  const filters = {};
  const gender = elements.gender.value;
  const stateCode = elements.stateInput.value.trim().toUpperCase();

  if (gender) {
    filters.gender = gender;
  }
  if (stateCode) {
    filters.state = stateCode;
  }
  if (elements.ageEnabled.checked) {
    filters.age_min = Number(elements.ageMin.value);
    filters.age_max = Number(elements.ageMax.value);
  }
  if (Object.keys(filters).length > 0) {
    request.filters = filters;
  }

  const includeFields = getSelectedFields();
  if (includeFields.length > 0) {
    request.include_fields = includeFields;
  }

  return request;
}

function getSelectedFields() {
  return Array.from(
    elements.fieldGroups.querySelectorAll('input[name="include_fields"]:checked'),
    (input) => input.value,
  );
}

function scheduleAutoGenerate() {
  if (!elements.autoGenerate.checked) {
    return;
  }

  const validationError = validateForm();
  if (validationError) {
    return;
  }

  window.clearTimeout(state.autoTimer);
  state.autoTimer = window.setTimeout(() => {
    generateProfiles("auto");
  }, 450);
}

async function generateProfiles(trigger) {
  const validationError = validateForm();
  if (validationError) {
    setStatus(validationError, "error");
    return;
  }

  const request = buildRequest();
  state.lastRequest = request;

  if (state.activeController) {
    state.activeController.abort();
  }

  const controller = new AbortController();
  state.activeController = controller;

  setStatus(statusCopyForTrigger(trigger), "loading");
  elements.generateButton.disabled = true;

  try {
    const response = await fetch("/v1/profiles/generate", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
      signal: controller.signal,
    });

    const payload = await response.json();
    if (!response.ok) {
      const message = payload?.error?.message || "请求失败";
      throw new Error(message);
    }

    state.lastResponse = payload;
    renderResponse(payload);

    const generatedAt = formatDateTime(payload?.meta?.generated_at);
    setStatus(`已生成 ${payload?.meta?.count || 0} 条资料，时间 ${generatedAt}。`, "success");
  } catch (error) {
    if (error.name === "AbortError") {
      return;
    }

    state.lastResponse = null;
    renderError(error.message);
    setStatus(error.message, "error");
  } finally {
    if (state.activeController === controller) {
      state.activeController = null;
    }
    elements.generateButton.disabled = Boolean(validateForm());
  }
}

function renderResponse(payload) {
  renderMeta(payload.meta);
  renderProfiles(payload.profiles || []);
  elements.responsePreview.textContent = JSON.stringify(payload, null, 2);
}

function renderMeta(meta) {
  if (!meta) {
    elements.metaGrid.innerHTML = "";
    return;
  }

  const cards = [
    ["Country", meta.country || "US"],
    ["Count", String(meta.count || 0)],
    ["Generated At", formatDateTime(meta.generated_at)],
  ];

  elements.metaGrid.innerHTML = cards
    .map(
      ([label, value]) => `
        <article class="meta-card">
          <span>${escapeHtml(label)}</span>
          <strong>${escapeHtml(value)}</strong>
        </article>
      `,
    )
    .join("");
}

function renderProfiles(profiles) {
  if (!profiles.length) {
    elements.profiles.innerHTML = `
      <article class="empty-state">
        <p>本次请求没有返回资料。</p>
      </article>
    `;
    return;
  }

  elements.profiles.innerHTML = profiles
    .map((profile, index) => renderProfileCard(profile, index))
    .join("");
}

function renderProfileCard(profile, index) {
  const initials = getInitials(profile.full_name);
  const subtitle = profile.email || profile.username || profile.occupation || "Synthetic Profile";
  const facts = [
    ["Age", profile.age],
    ["Birth Date", profile.birth_date],
    ["Occupation", profile.occupation],
    ["Username", profile.username],
  ]
    .filter(([, value]) => value !== undefined && value !== null && value !== "")
    .map(
      ([label, value]) => `
        <article class="fact">
          <span>${escapeHtml(label)}</span>
          <strong>${escapeHtml(String(value))}</strong>
        </article>
      `,
    )
    .join("");

  const detailBlocks = [
    renderDetailBlock("Contact", [
      ["Email", profile.email],
      ["Phone", profile.phone?.national],
      ["E.164", profile.phone?.e164],
      ["Area Code", profile.phone?.area_code],
      ["Time Zone", profile.time_zone],
    ]),
    renderDetailBlock("Address", [
      ["Street 1", profile.address?.street_line1],
      ["Street 2", profile.address?.street_line2],
      ["City", profile.address?.city],
      ["State", profile.address?.state],
      ["State Code", profile.address?.state_code],
      ["Postal Code", profile.address?.postal_code],
    ]),
    renderDetailBlock("Employment", [
      ["Status", profile.employment?.status],
      ["Company", profile.employment?.company_name],
      ["Domain", profile.employment?.company_domain],
      ["Title", profile.employment?.job_title],
      ["Department", profile.employment?.department],
      ["Work Email", profile.employment?.work_email],
      ["Industry", profile.employment?.industry],
      ["Seniority", profile.employment?.seniority],
      ["Start Date", profile.employment?.start_date],
      ["Work Mode", profile.employment?.work_mode],
      ["Manager", profile.employment?.manager_name],
      ["Office City", profile.employment?.office_city],
      ["Office State", profile.employment?.office_state],
    ]),
    renderDetailBlock("Education", [
      ["Level", profile.education?.level],
      ["Degree Type", profile.education?.degree_type],
      ["School Name", profile.education?.school_name],
      ["Major", profile.education?.major],
      ["Graduation Year", profile.education?.graduation_year],
      ["Student Status", profile.education?.student_status],
      ["School City", profile.education?.school_city],
      ["School State", profile.education?.school_state],
      ["School Email", profile.education?.school_email],
    ]),
    renderDetailBlock("Sensitive", [
      ["SSN", profile.sensitive?.ssn],
      ["DL Number", profile.sensitive?.driver_license?.number],
      ["DL State", profile.sensitive?.driver_license?.state_code],
      ["Passport No.", profile.sensitive?.passport?.number],
      ["Passport Country", profile.sensitive?.passport?.country],
      ["Card Brand", profile.sensitive?.bank_card?.brand],
      ["Card Number", profile.sensitive?.bank_card?.number],
      ["Expiry Month", profile.sensitive?.bank_card?.expiry_month],
      ["Expiry Year", profile.sensitive?.bank_card?.expiry_year],
      ["CVV", profile.sensitive?.bank_card?.cvv],
    ]),
  ]
    .filter(Boolean)
    .join("");

  const tags = [
    profile.gender,
    profile.marital_status,
    profile.country,
    profile.address?.state_code,
  ]
    .filter(Boolean)
    .map((value) => `<span class="tag">${escapeHtml(String(value))}</span>`)
    .join("");

  return `
    <article class="profile-card">
      <header class="profile-head">
        <div class="identity-block">
          <div class="identity-mark" aria-hidden="true">${escapeHtml(initials)}</div>
          <div>
            <p class="profile-index">Profile ${String(index + 1).padStart(2, "0")}</p>
            <h3>${escapeHtml(profile.full_name || "Synthetic Profile")}</h3>
            <p class="profile-subtitle">${escapeHtml(String(subtitle))}</p>
          </div>
        </div>
        <div class="profile-tags">${tags}</div>
      </header>

      <section class="facts-grid">${facts}</section>
      <section class="detail-grid">${detailBlocks}</section>
    </article>
  `;
}

function renderDetailBlock(title, entries) {
  const rows = entries
    .filter(([, value]) => value !== undefined && value !== null && value !== "")
    .map(
      ([label, value]) => `
        <div>
          <dt>${escapeHtml(String(label))}</dt>
          <dd>${escapeHtml(String(value))}</dd>
        </div>
      `,
    )
    .join("");

  if (!rows) {
    return "";
  }

  return `
    <section class="detail-block">
      <h4>${escapeHtml(title)}</h4>
      <dl>${rows}</dl>
    </section>
  `;
}

function renderError(message) {
  elements.metaGrid.innerHTML = "";
  elements.profiles.innerHTML = `
    <article class="empty-state">
      <p>${escapeHtml(message)}</p>
    </article>
  `;
  elements.responsePreview.textContent = JSON.stringify(
    {
      error: message,
      request: state.lastRequest,
    },
    null,
    2,
  );
}

function resetForm() {
  elements.count.value = "3";
  elements.gender.value = "random";
  elements.stateInput.value = "";
  elements.ageEnabled.checked = false;
  elements.ageMin.value = "24";
  elements.ageMax.value = "42";
  elements.autoGenerate.checked = true;

  const inputs = elements.fieldGroups.querySelectorAll('input[name="include_fields"]');
  inputs.forEach((input) => {
    input.checked = DEFAULT_FIELDS.has(input.value);
  });
}

function statusCopyForTrigger(trigger) {
  if (trigger === "initial") {
    return "页面已加载，正在自动生成示例资料。";
  }
  if (trigger === "auto") {
    return "检测到参数变动，正在自动刷新结果。";
  }
  if (trigger === "reset") {
    return "已恢复默认参数，正在重新生成。";
  }

  return "正在生成合成资料。";
}

function setStatus(message, tone) {
  elements.statusBanner.className = `status-banner ${tone}`;
  elements.statusBanner.textContent = message;
}

function setFormMessage(message, tone) {
  elements.formMessage.className = tone ? `form-message ${tone}` : "form-message";
  elements.formMessage.textContent = message;
}

async function copyText(text, successMessage) {
  try {
    await navigator.clipboard.writeText(text);
    setStatus(successMessage, "success");
  } catch (error) {
    setStatus(`复制失败: ${error.message}`, "error");
  }
}

function downloadJSON(payload) {
  const blob = new Blob([JSON.stringify(payload, null, 2)], {
    type: "application/json;charset=utf-8",
  });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `synthetic-profiles-${Date.now()}.json`;
  link.click();
  URL.revokeObjectURL(url);
  setStatus("响应 JSON 已下载。", "success");
}

function formatDateTime(value) {
  if (!value) {
    return "N/A";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(date);
}

function getInitials(value) {
  const parts = String(value || "")
    .trim()
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2);

  if (!parts.length) {
    return "FP";
  }

  return parts.map((part) => part[0].toUpperCase()).join("");
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
