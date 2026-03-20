const STATES = [
  ["AL", "Alabama"],
  ["AK", "Alaska"],
  ["AZ", "Arizona"],
  ["AR", "Arkansas"],
  ["CA", "California"],
  ["CO", "Colorado"],
  ["CT", "Connecticut"],
  ["DE", "Delaware"],
  ["FL", "Florida"],
  ["GA", "Georgia"],
  ["HI", "Hawaii"],
  ["ID", "Idaho"],
  ["IL", "Illinois"],
  ["IN", "Indiana"],
  ["IA", "Iowa"],
  ["KS", "Kansas"],
  ["KY", "Kentucky"],
  ["LA", "Louisiana"],
  ["ME", "Maine"],
  ["MD", "Maryland"],
  ["MA", "Massachusetts"],
  ["MI", "Michigan"],
  ["MN", "Minnesota"],
  ["MS", "Mississippi"],
  ["MO", "Missouri"],
  ["MT", "Montana"],
  ["NE", "Nebraska"],
  ["NV", "Nevada"],
  ["NH", "New Hampshire"],
  ["NJ", "New Jersey"],
  ["NM", "New Mexico"],
  ["NY", "New York"],
  ["NC", "North Carolina"],
  ["ND", "North Dakota"],
  ["OH", "Ohio"],
  ["OK", "Oklahoma"],
  ["OR", "Oregon"],
  ["PA", "Pennsylvania"],
  ["RI", "Rhode Island"],
  ["SC", "South Carolina"],
  ["SD", "South Dakota"],
  ["TN", "Tennessee"],
  ["TX", "Texas"],
  ["UT", "Utah"],
  ["VT", "Vermont"],
  ["VA", "Virginia"],
  ["WA", "Washington"],
  ["WV", "West Virginia"],
  ["WI", "Wisconsin"],
  ["WY", "Wyoming"],
  ["DC", "District of Columbia"],
];

const FIELD_GROUPS = [
  {
    title: "Education",
    description: "Academic details that make profiles feel production-like.",
    items: [
      ["education.level", "Level", "Bachelor, master, associate, or similar"],
      ["education.degree_type", "Degree Type", "B.S., M.S., MBA, and related labels"],
      ["education.school_name", "School Name", "University or college name"],
      ["education.major", "Major", "Field of study"],
      ["education.graduation_year", "Graduation Year", "Graduation timeline"],
      ["education.student_status", "Student Status", "Alumni or currently enrolled"],
      ["education.school_city", "School City", "Campus city"],
      ["education.school_state", "School State", "Campus state"],
      ["education.school_email", "School Email", "Education-linked email address"],
    ],
  },
  {
    title: "Employment",
    description: "Enhanced work metadata for internal tools and test systems.",
    items: [
      ["employment.industry", "Industry", "Business sector"],
      ["employment.seniority", "Seniority", "IC, manager, director, and more"],
      ["employment.start_date", "Start Date", "Join date"],
      ["employment.work_mode", "Work Mode", "Remote, hybrid, or on-site"],
      ["employment.manager_name", "Manager Name", "Reporting line"],
      ["employment.office_city", "Office City", "Office city"],
      ["employment.office_state", "Office State", "Office state"],
    ],
  },
  {
    title: "Sensitive",
    description: "Synthetic sensitive fields for isolated QA and staging only.",
    items: [
      ["sensitive.ssn", "SSN", "Social security number"],
      ["sensitive.driver_license.number", "Driver License No.", "License identifier"],
      ["sensitive.driver_license.state_code", "DL State", "Issuing state"],
      ["sensitive.passport.number", "Passport No.", "Passport identifier"],
      ["sensitive.passport.country", "Passport Country", "Issuing country"],
      ["sensitive.bank_card.brand", "Card Brand", "Payment network"],
      ["sensitive.bank_card.number", "Card Number", "Synthetic PAN"],
      ["sensitive.bank_card.expiry_month", "Expiry Month", "Month component"],
      ["sensitive.bank_card.expiry_year", "Expiry Year", "Year component"],
      ["sensitive.bank_card.cvv", "CVV", "Card verification value"],
    ],
  },
];

const DEFAULT_FIELDS = new Set([
  "education.degree_type",
  "education.school_name",
  "education.major",
  "education.graduation_year",
  "employment.industry",
  "employment.seniority",
  "employment.start_date",
  "employment.work_mode",
  "sensitive.ssn",
]);

const state = {
  activeController: null,
  autoTimer: 0,
  lastLatency: null,
  lastRequest: null,
  lastResponse: null,
  lastResponseText: "",
  sessionId: buildSessionId(),
};

const elements = {};

document.addEventListener("DOMContentLoaded", () => {
  cacheElements();
  populateStateOptions();
  renderFieldGroups();
  bindEvents();
  initializeSession();
  toggleAgeInputs();
  syncUI();
  generateProfiles("initial");
});

function cacheElements() {
  elements.form = document.getElementById("generator-form");
  elements.count = document.getElementById("count");
  elements.countValue = document.getElementById("count-value");
  elements.gender = document.getElementById("gender");
  elements.stateInput = document.getElementById("state");
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
  elements.advancedPanel = document.getElementById("advanced-panel");
  elements.advancedSummaryText = document.getElementById("advanced-summary-text");
  elements.fieldGroups = document.getElementById("field-groups");
  elements.fieldCount = document.getElementById("field-count");
  elements.requestPreview = document.getElementById("request-preview");
  elements.responsePreview = document.getElementById("response-preview");
  elements.statusBanner = document.getElementById("status-banner");
  elements.metaGrid = document.getElementById("meta-grid");
  elements.profiles = document.getElementById("profiles");
  elements.nextBatchButton = document.getElementById("next-batch-button");
  elements.latencyValue = document.getElementById("latency-pill-value");
  elements.generatedTotalValue = document.getElementById("generated-total-value");
  elements.ageSummaryChip = document.getElementById("age-summary-chip");
  elements.fieldSummaryChip = document.getElementById("field-summary-chip");
  elements.apiUsage = document.getElementById("api-usage");
  elements.sessionId = document.getElementById("session-id");
  elements.responseDisclosure = document.getElementById("response-disclosure");
  elements.navItems = Array.from(document.querySelectorAll(".nav-item"));
}

function populateStateOptions() {
  const options = STATES.map(
    ([code, name]) => `<option value="${code}">${escapeHtml(name)}</option>`,
  ).join("");
  elements.stateInput.insertAdjacentHTML("beforeend", options);
}

function renderFieldGroups() {
  elements.fieldGroups.innerHTML = FIELD_GROUPS.map((group) => {
    const items = group.items
      .map(([value, label, hint]) => {
        const checked = DEFAULT_FIELDS.has(value) ? "checked" : "";
        return `
          <label class="checkbox">
            <input type="checkbox" name="include_fields" value="${escapeHtml(value)}" ${checked} />
            <span class="checkbox-copy">
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

    await copyText(JSON.stringify(buildRequest(), null, 2), "Request JSON copied.");
  });

  elements.copyResponseButton.addEventListener("click", async () => {
    if (!state.lastResponse) {
      setStatus("No response is available to copy yet.", "idle");
      return;
    }

    await copyText(JSON.stringify(state.lastResponse, null, 2), "Response JSON copied.");
  });

  elements.downloadResponseButton.addEventListener("click", () => {
    if (!state.lastResponse) {
      setStatus("No response is available to export yet.", "idle");
      return;
    }

    downloadJSON(state.lastResponse);
  });

  elements.nextBatchButton.addEventListener("click", async () => {
    await generateProfiles("manual");
  });

  elements.profiles.addEventListener("click", async (event) => {
    const target = event.target.closest("[data-action]");
    if (!target) {
      return;
    }

    const action = target.dataset.action;
    if (action === "toggle-sensitive") {
      toggleSensitiveValue(target);
      return;
    }

    if (action === "copy-profile") {
      await copyText(target.dataset.value || "", "Profile reference copied.");
    }
  });

  elements.navItems.forEach((item) => {
    item.addEventListener("click", () => {
      openPanel(item.dataset.navTarget || "");
      setActiveNavItem(item.dataset.navTarget || "");
    });
  });

  elements.responseDisclosure.addEventListener("toggle", () => {
    if (elements.responseDisclosure.open) {
      syncResponsePreview();
    }
  });
}

function initializeSession() {
  elements.sessionId.textContent = state.sessionId;
}

function syncUI() {
  toggleAgeInputs();
  updateCountBadge();
  updateAdvancedSummary();
  updateDashboardSummary();
  updateRequestPreview();
}

function toggleAgeInputs() {
  const enabled = elements.ageEnabled.checked;
  elements.ageMin.disabled = !enabled;
  elements.ageMax.disabled = !enabled;
}

function updateCountBadge() {
  elements.countValue.textContent = elements.count.value;
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

function updateDashboardSummary() {
  elements.ageSummaryChip.textContent = elements.ageEnabled.checked
    ? `${elements.ageMin.value} - ${elements.ageMax.value}`
    : "All Ages";

  const fieldSummary = getFieldSummaryLabels();
  elements.fieldSummaryChip.textContent = fieldSummary.join(", ");

  const usage = Math.min(97.4, 8 + Number(elements.count.value) * 0.248);
  elements.apiUsage.textContent = `${usage.toFixed(1)}%`;
}

function updateRequestPreview() {
  const validationError = validateForm();
  if (validationError) {
    elements.requestPreview.textContent = validationError;
    elements.generateButton.disabled = true;
    setFormMessage(validationError, "error");
    if (validationError.includes("Age")) {
      elements.advancedPanel.open = true;
    }
    return;
  }

  const request = buildRequest();
  state.lastRequest = request;
  elements.requestPreview.textContent = JSON.stringify(request, null, 2);
  elements.generateButton.disabled = false;
  setFormMessage("Ready to generate a fresh synthetic batch.", "success");
}

function validateForm() {
  const count = Number(elements.count.value);
  if (!Number.isInteger(count) || count < 1 || count > 20) {
    return "Result Count must be an integer between 1 and 20.";
  }

  if (elements.stateInput.value && !isKnownStateCode(elements.stateInput.value)) {
    return "State must be a supported US state.";
  }

  if (elements.ageEnabled.checked) {
    const ageMin = Number(elements.ageMin.value);
    const ageMax = Number(elements.ageMax.value);

    if (!Number.isInteger(ageMin) || !Number.isInteger(ageMax)) {
      return "Age Range requires both a minimum and maximum value.";
    }
    if (ageMin < 0 || ageMax < 0 || ageMin > 100 || ageMax > 100) {
      return "Age Range must stay within 0 and 100.";
    }
    if (ageMin > ageMax) {
      return "Age minimum cannot be greater than age maximum.";
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
  if (elements.gender.value && elements.gender.value !== "random") {
    filters.gender = elements.gender.value;
  }
  if (elements.stateInput.value) {
    filters.state = elements.stateInput.value;
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

function getFieldSummaryLabels() {
  const selected = getSelectedFields();
  const labels = [];
  const hasSensitiveFields = selected.some((value) => value.startsWith("sensitive."));
  const hasExpandedSensitiveFields = selected.some(
    (value) => value.startsWith("sensitive.") && value !== "sensitive.ssn",
  );

  if (selected.some((value) => value.startsWith("employment."))) {
    labels.push("Employment");
  }
  if (selected.some((value) => value.startsWith("education."))) {
    labels.push("Education");
  }
  if (hasExpandedSensitiveFields) {
    labels.push("Sensitive");
  } else if (selected.includes("sensitive.ssn")) {
    labels.push("SSN");
  } else if (hasSensitiveFields) {
    labels.push("Sensitive");
  }

  labels.push("Contact");

  return labels.slice(0, 4);
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
  const startedAt = performance.now();

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
      const message = payload?.error?.message || "Request failed.";
      throw new Error(message);
    }

    state.lastLatency = Math.max(1, Math.round(performance.now() - startedAt));
    state.lastResponse = payload;
    state.lastResponseText = "";
    renderResponse(payload);
    updateTopMetrics(payload?.meta);

    const generatedAt = formatDateTime(payload?.meta?.generated_at);
    setStatus(
      `Generated ${payload?.meta?.count || 0} profiles at ${generatedAt}.`,
      "success",
    );
  } catch (error) {
    if (error.name === "AbortError") {
      return;
    }

    state.lastResponse = null;
    state.lastResponseText = "";
    renderError(error.message);
    updateTopMetrics();
    setStatus(error.message, "error");
  } finally {
    if (state.activeController === controller) {
      state.activeController = null;
    }
    elements.generateButton.disabled = Boolean(validateForm());
  }
}

function updateTopMetrics(meta) {
  elements.latencyValue.textContent = state.lastLatency
    ? `${state.lastLatency} ms latency`
    : "-- ms latency";

  const total = meta?.count || 0;
  elements.generatedTotalValue.textContent = `${total} ${
    total === 1 ? "Profile" : "Profiles"
  } Generated`;
}

function renderResponse(payload) {
  renderMeta(payload.meta);
  renderProfiles(payload.profiles || []);
  syncResponsePreview();
}

function syncResponsePreview() {
  if (!state.lastResponse && !state.lastResponseText) {
    elements.responsePreview.textContent = "";
    return;
  }

  if (!elements.responseDisclosure.open) {
    elements.responsePreview.textContent = "Expand to inspect the raw JSON response.";
    return;
  }

  if (!state.lastResponseText && state.lastResponse) {
    state.lastResponseText = JSON.stringify(state.lastResponse, null, 2);
  }

  elements.responsePreview.textContent = state.lastResponseText;
}

function renderMeta(meta) {
  if (!meta) {
    elements.metaGrid.innerHTML = "";
    return;
  }

  const cards = [
    ["Country", meta.country || "US"],
    ["Generated At", formatDateTime(meta.generated_at)],
    [
      "Filters",
      `${humanizeGender(elements.gender.value)} / ${
        elements.stateInput.value ? humanizeState(elements.stateInput.value) : "All States"
      }`,
    ],
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
        <p>No profiles were returned for this request.</p>
      </article>
    `;
    return;
  }

  elements.profiles.innerHTML = profiles
    .map((profile, index) => renderProfileCard(profile, index))
    .join("");
}

function renderProfileCard(profile, index) {
  const subtitle = [
    humanizeGender(profile.gender),
    `${profile.age || "N/A"} yrs`,
    `${profile.address?.state || "Unknown"}, ${countryName(profile.country)}`,
  ].join(" • ");

  const employmentHeadline = profile.employment?.job_title || profile.occupation || "Unassigned role";
  const employmentMeta = buildEmploymentMeta(profile);
  const educationHeadline = buildEducationHeadline(profile.education);
  const educationMeta = buildEducationMeta(profile.education);
  const hasExpandedSensitiveData = hasExtendedSensitiveDetails(profile.sensitive);
  const sensitiveHeadline = buildSensitiveHeadline(profile.sensitive);
  const sensitiveMeta = buildSensitiveMeta(profile.sensitive);
  const maskedSSN = profile.sensitive?.ssn
    ? maskSensitive(profile.sensitive.ssn, "ssn")
    : "Not requested";
  const hasPremiumData = Boolean(
    hasSensitiveDetails(profile.sensitive) ||
      profile.education ||
      profile.employment?.industry ||
      profile.employment?.seniority ||
      profile.employment?.start_date,
  );

  const addressSecondary = [
    formatCityStatePostal(profile.address),
    countryName(profile.country),
  ]
    .filter(Boolean)
    .join("<br />");

  return `
    <article class="profile-card">
      <header class="profile-card-header">
        <div class="profile-title">
          <div class="profile-title-row">
            <h3>${escapeHtml(profile.full_name || `Profile ${index + 1}`)}</h3>
            <div class="profile-badges">
              ${
                hasPremiumData
                  ? '<span class="badge premium">Premium</span>'
                  : '<span class="badge premium">Base</span>'
              }
              <span class="badge verified">Verified</span>
            </div>
          </div>
          <p class="profile-subline">${escapeHtml(subtitle)}</p>
        </div>

        <button class="card-menu" type="button" aria-label="Card actions" disabled>
          ⋮
        </button>
      </header>

      <div class="profile-card-body">
        <section class="info-stack">
          ${renderInfoRow("mail", profile.email || "No email", profile.employment?.work_email || "")}
          ${renderInfoRow("phone", profile.phone?.national || "No phone", profile.phone?.e164 || "")}
          ${renderInfoRow("pin", formatStreetAddress(profile.address), addressSecondary, true)}
        </section>

        <section class="section-stack">
          ${renderSectionBlock("Employment", employmentHeadline, employmentMeta)}
          ${renderSectionBlock("Education", educationHeadline, educationMeta)}
          ${
            hasExpandedSensitiveData
              ? renderSectionBlock("Sensitive", sensitiveHeadline, sensitiveMeta)
              : ""
          }
        </section>
      </div>

      <footer class="profile-card-footer">
        <div class="footer-field" data-sensitive-container>
          <span class="footer-label">SSN</span>
          <span class="footer-value ${profile.sensitive?.ssn ? "" : "subtle"}" data-sensitive-value>
            ${escapeHtml(maskedSSN)}
          </span>
          ${
            profile.sensitive?.ssn
              ? `
                <button
                  class="link-button"
                  type="button"
                  data-action="toggle-sensitive"
                  data-visible="false"
                  data-value="${escapeHtml(profile.sensitive.ssn)}"
                  data-mask="${escapeHtml(maskedSSN)}"
                >
                  Reveal
                </button>
              `
              : ""
          }
        </div>

        <div class="footer-field">
          <span class="footer-label">DOB</span>
          <span class="footer-value">${escapeHtml(profile.birth_date || "N/A")}</span>
        </div>

        <button
          class="icon-button"
          type="button"
          data-action="copy-profile"
          data-value="${escapeHtml(profile.id || profile.email || profile.full_name || "")}"
          aria-label="Copy profile reference"
        >
          ⧉
        </button>
      </footer>
    </article>
  `;
}

function renderInfoRow(icon, primary, secondary, secondaryIsHtml = false) {
  if (!primary && !secondary) {
    return "";
  }

  const secondaryMarkup = secondary
    ? secondaryIsHtml
      ? secondary
      : escapeHtml(secondary)
    : "";

  return `
    <div class="info-row">
      <span class="info-icon" aria-hidden="true">${renderIcon(icon)}</span>
      <div class="info-text">
        <div class="info-primary">${escapeHtml(primary || "N/A")}</div>
        ${secondaryMarkup ? `<div class="info-secondary">${secondaryMarkup}</div>` : ""}
      </div>
    </div>
  `;
}

function renderSectionBlock(title, headline, meta) {
  return `
    <section class="section-block">
      <p class="eyebrow">${escapeHtml(title)}</p>
      <strong>${escapeHtml(headline || "Not available")}</strong>
      <div class="section-copy">${meta || "Not requested for this batch."}</div>
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
  state.lastResponseText = JSON.stringify(
    {
      error: message,
      request: state.lastRequest,
    },
    null,
    2,
  );
  elements.responsePreview.textContent = elements.responseDisclosure.open
    ? state.lastResponseText
    : "Expand to inspect the error payload.";
}

function resetForm() {
  elements.count.value = "4";
  elements.gender.value = "random";
  elements.stateInput.value = "";
  elements.ageEnabled.checked = true;
  elements.ageMin.value = "18";
  elements.ageMax.value = "65";
  elements.autoGenerate.checked = true;

  const inputs = elements.fieldGroups.querySelectorAll('input[name="include_fields"]');
  inputs.forEach((input) => {
    input.checked = DEFAULT_FIELDS.has(input.value);
  });
}

function statusCopyForTrigger(trigger) {
  if (trigger === "initial") {
    return "Initializing the workbench and generating a sample batch.";
  }
  if (trigger === "auto") {
    return "Filters changed. Refreshing the result set.";
  }
  if (trigger === "reset") {
    return "Defaults restored. Generating a clean batch.";
  }

  return "Generating synthetic profiles.";
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
    setStatus(`Copy failed: ${error.message}`, "error");
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
  setStatus("Response JSON exported.", "success");
}

function toggleSensitiveValue(button) {
  const container = button.closest("[data-sensitive-container]");
  const valueNode = container?.querySelector("[data-sensitive-value]");
  if (!container || !valueNode) {
    return;
  }

  const rawValue = button.dataset.value || "";
  const maskedValue = button.dataset.mask || maskSensitive(rawValue);
  const visible = button.dataset.visible === "true";
  valueNode.textContent = visible ? maskedValue : rawValue;
  button.dataset.visible = visible ? "false" : "true";
  button.textContent = visible ? "Reveal" : "Hide";
}

function openPanel(targetId) {
  const panel = document.getElementById(targetId);
  if (!panel) {
    return;
  }

  if (panel instanceof HTMLDetailsElement) {
    panel.open = true;
  }

  panel.scrollIntoView({
    behavior: prefersReducedMotion() ? "auto" : "smooth",
    block: "start",
  });
}

function setActiveNavItem(targetId) {
  elements.navItems.forEach((item) => {
    item.classList.toggle("active", item.dataset.navTarget === targetId);
  });
}

function buildEmploymentMeta(profile) {
  const parts = [];
  if (profile.employment?.company_name) {
    parts.push(profile.employment.company_name);
  }

  const experience = getExperienceYears(profile.employment?.start_date);
  if (experience) {
    parts.push(`${experience} Years Exp.`);
  } else if (profile.employment?.seniority) {
    parts.push(titleCase(profile.employment.seniority));
  }

  if (profile.employment?.industry) {
    parts.push(profile.employment.industry);
  } else if (profile.employment?.work_mode) {
    parts.push(titleCase(profile.employment.work_mode));
  }

  return parts.length ? escapeHtml(parts.join(" • ")) : "Not requested for this batch.";
}

function buildEducationHeadline(education) {
  if (!education) {
    return "Education hidden";
  }

  const parts = [];
  if (education.degree_type) {
    parts.push(education.degree_type);
  } else if (education.level) {
    parts.push(titleCase(education.level));
  }
  if (education.major) {
    parts.push(education.major);
  }

  if (parts.length) {
    return parts.join(" ");
  }
  if (education.school_name) {
    return education.school_name;
  }

  return "Education included";
}

function buildEducationMeta(education) {
  if (!education) {
    return "Not requested for this batch.";
  }

  const parts = [];
  if (education.school_name) {
    parts.push(education.school_name);
  }
  if (education.graduation_year) {
    parts.push(String(education.graduation_year));
  }
  if (education.student_status) {
    parts.push(titleCase(education.student_status));
  }

  return parts.length ? escapeHtml(parts.join(", ")) : "Education included without extra details.";
}

function buildSensitiveHeadline(sensitive) {
  if (!hasExtendedSensitiveDetails(sensitive)) {
    return "Sensitive fields hidden";
  }

  const labels = [];
  if (hasBankCardDetails(sensitive?.bank_card)) {
    labels.push(sensitive.bank_card?.brand ? `${sensitive.bank_card.brand} card` : "Bank card");
  }
  if (hasDriverLicenseDetails(sensitive?.driver_license)) {
    labels.push("Driver license");
  }
  if (hasPassportDetails(sensitive?.passport)) {
    labels.push("Passport");
  }

  if (labels.length === 1) {
    return labels[0];
  }
  if (labels.length === 2) {
    return `${labels[0]} and ${labels[1]}`;
  }

  return `${labels[0]}, ${labels[1]}, and ${labels[2]}`;
}

function buildSensitiveMeta(sensitive) {
  if (!hasExtendedSensitiveDetails(sensitive)) {
    return "Not requested for this batch.";
  }

  const groups = [];

  if (hasBankCardDetails(sensitive?.bank_card)) {
    const details = [];
    details.push(renderSensitiveDetail("Brand", sensitive.bank_card?.brand || ""));
    details.push(
      renderSensitiveDetail("Number", sensitive.bank_card?.number || "", {
        revealable: true,
        maskedValue: maskSensitive(sensitive.bank_card?.number || "", "card"),
      }),
    );

    const expiry = formatCardExpiry(sensitive.bank_card);
    details.push(renderSensitiveDetail("Expiry", expiry));
    details.push(
      renderSensitiveDetail("CVV", sensitive.bank_card?.cvv || "", {
        revealable: true,
        maskedValue: maskSensitive(sensitive.bank_card?.cvv || "", "full"),
      }),
    );

    groups.push(renderSensitiveGroup("Bank Card", details.join("")));
  }

  if (hasDriverLicenseDetails(sensitive?.driver_license)) {
    const details = [];
    details.push(
      renderSensitiveDetail("Number", sensitive.driver_license?.number || "", {
        revealable: true,
      }),
    );
    details.push(renderSensitiveDetail("State", sensitive.driver_license?.state_code || ""));
    groups.push(renderSensitiveGroup("Driver License", details.join("")));
  }

  if (hasPassportDetails(sensitive?.passport)) {
    const details = [];
    details.push(
      renderSensitiveDetail("Number", sensitive.passport?.number || "", {
        revealable: true,
      }),
    );
    details.push(renderSensitiveDetail("Country", sensitive.passport?.country || ""));
    groups.push(renderSensitiveGroup("Passport", details.join("")));
  }

  return groups.length
    ? `<div class="sensitive-stack">${groups.join("")}</div>`
    : "Not requested for this batch.";
}

function renderSensitiveGroup(label, details) {
  if (!details) {
    return "";
  }

  return `
    <div class="sensitive-group">
      <span class="sensitive-group-label">${escapeHtml(label)}</span>
      <div class="sensitive-group-body">${details}</div>
    </div>
  `;
}

function renderSensitiveDetail(label, value, options = {}) {
  if (!value) {
    return "";
  }

  const safeLabel = escapeHtml(label);
  const safeValue = escapeHtml(value);

  if (!options.revealable) {
    return `
      <div class="sensitive-detail">
        <span class="sensitive-key">${safeLabel}</span>
        <span class="sensitive-value">${safeValue}</span>
      </div>
    `;
  }

  const maskedValue = escapeHtml(options.maskedValue || maskSensitive(value));
  return `
    <div class="sensitive-detail" data-sensitive-container>
      <span class="sensitive-key">${safeLabel}</span>
      <span class="sensitive-value" data-sensitive-value>${maskedValue}</span>
      <button
        class="link-button"
        type="button"
        data-action="toggle-sensitive"
        data-visible="false"
        data-value="${safeValue}"
        data-mask="${maskedValue}"
      >
        Reveal
      </button>
    </div>
  `;
}

function formatCardExpiry(bankCard) {
  if (!bankCard) {
    return "";
  }

  const month = bankCard.expiry_month;
  const year = bankCard.expiry_year;
  if (!month && !year) {
    return "";
  }
  if (!month) {
    return String(year);
  }
  if (!year) {
    return String(month).padStart(2, "0");
  }

  return `${String(month).padStart(2, "0")}/${year}`;
}

function hasSensitiveDetails(sensitive) {
  return Boolean(
    sensitive?.ssn ||
      hasBankCardDetails(sensitive?.bank_card) ||
      hasDriverLicenseDetails(sensitive?.driver_license) ||
      hasPassportDetails(sensitive?.passport),
  );
}

function hasExtendedSensitiveDetails(sensitive) {
  return Boolean(
    hasBankCardDetails(sensitive?.bank_card) ||
      hasDriverLicenseDetails(sensitive?.driver_license) ||
      hasPassportDetails(sensitive?.passport),
  );
}

function hasBankCardDetails(bankCard) {
  return Boolean(
    bankCard?.brand ||
      bankCard?.number ||
      bankCard?.expiry_month ||
      bankCard?.expiry_year ||
      bankCard?.cvv,
  );
}

function hasDriverLicenseDetails(driverLicense) {
  return Boolean(driverLicense?.number || driverLicense?.state_code);
}

function hasPassportDetails(passport) {
  return Boolean(passport?.number || passport?.country);
}

function getExperienceYears(startDate) {
  if (!startDate) {
    return "";
  }

  const start = new Date(startDate);
  if (Number.isNaN(start.getTime())) {
    return "";
  }

  const now = new Date();
  let years = now.getFullYear() - start.getFullYear();
  const monthDelta = now.getMonth() - start.getMonth();
  if (monthDelta < 0 || (monthDelta === 0 && now.getDate() < start.getDate())) {
    years -= 1;
  }

  return years > 0 ? String(years) : "1";
}

function formatStreetAddress(address) {
  if (!address) {
    return "No address";
  }

  return [address.street_line1, address.street_line2].filter(Boolean).join(", ");
}

function formatCityStatePostal(address) {
  if (!address) {
    return "";
  }

  return [address.city, address.state_code || address.state, address.postal_code]
    .filter(Boolean)
    .join(", ");
}

function maskSensitive(value, mode = "tail") {
  const normalized = String(value || "");
  if (!normalized) {
    return "";
  }

  if (mode === "full") {
    return "*".repeat(normalized.length);
  }
  if (mode === "card") {
    return normalized.length > 4 ? `**** ${normalized.slice(-4)}` : "*".repeat(normalized.length);
  }

  const digits = normalized.replace(/\D/g, "");
  if (mode === "ssn" || digits.length === 9) {
    return `***-**-${digits.slice(-4)}`;
  }
  if (normalized.length <= 4) {
    return "*".repeat(normalized.length);
  }

  return `${"*".repeat(Math.max(0, normalized.length - 4))}${normalized.slice(-4)}`;
}

function isKnownStateCode(value) {
  return STATES.some(([code]) => code === value);
}

function humanizeState(value) {
  return STATES.find(([code]) => code === value)?.[1] || value;
}

function humanizeGender(value) {
  if (value === "female") {
    return "Female";
  }
  if (value === "male") {
    return "Male";
  }
  return "All";
}

function countryName(value) {
  if (value === "US") {
    return "United States";
  }
  return value || "Unknown";
}

function formatDateTime(value) {
  if (!value) {
    return "N/A";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function buildSessionId() {
  return `${Math.floor(Math.random() * 9000 + 1000)}-${Math.floor(
    Math.random() * 90 + 10,
  )}-${Math.floor(Math.random() * 90 + 10)}`;
}

function renderIcon(name) {
  if (name === "mail") {
    return `
      <svg viewBox="0 0 16 16" focusable="false" aria-hidden="true">
        <path d="M2 3h12v10H2z" fill="none" stroke="currentColor" stroke-width="1.3"/>
        <path d="m2.5 4 5.5 4 5.5-4" fill="none" stroke="currentColor" stroke-width="1.3"/>
      </svg>
    `;
  }
  if (name === "phone") {
    return `
      <svg viewBox="0 0 16 16" focusable="false" aria-hidden="true">
        <path
          d="M5 2.8c.2-.4.7-.6 1.2-.4l1.5.6c.4.2.7.7.6 1.1l-.4 1.7a1 1 0 0 0 .2.9l1.1 1.1a1 1 0 0 0 .9.2l1.7-.4c.4-.1.9.1 1.1.6l.6 1.5c.2.5 0 1-.4 1.2l-1 .5c-.9.4-1.9.3-2.7-.2A15.4 15.4 0 0 1 3.7 7c-.5-.8-.6-1.8-.2-2.7z"
          fill="none"
          stroke="currentColor"
          stroke-width="1.3"
          stroke-linejoin="round"
        />
      </svg>
    `;
  }
  if (name === "pin") {
    return `
      <svg viewBox="0 0 16 16" focusable="false" aria-hidden="true">
        <path
          d="M8 14s4-4.3 4-7.2A4 4 0 1 0 4 6.8C4 9.7 8 14 8 14Z"
          fill="none"
          stroke="currentColor"
          stroke-width="1.3"
        />
        <circle cx="8" cy="6.4" r="1.4" fill="none" stroke="currentColor" stroke-width="1.3" />
      </svg>
    `;
  }

  return "";
}

function titleCase(value) {
  return String(value || "")
    .split(/[_\s-]+/)
    .filter(Boolean)
    .map((part) => part[0].toUpperCase() + part.slice(1))
    .join(" ");
}

function prefersReducedMotion() {
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
