document.addEventListener("DOMContentLoaded", function () {
  const elements = {
    form: document.getElementById("uploadForm"),
    progress: document.getElementById("progress"),
    result: document.getElementById("result"),
    formatSelect: document.getElementById("formatSelect"),
    qualitySelect: document.getElementById("qualitySelect"),
    widthInput: document.getElementById("widthInput"),
    heightInput: document.getElementById("heightInput"),
    removeBackground: document.getElementById("removeBackground"),
    bgRemovalOptions: document.getElementById("bgRemovalOptions"),
    resizeModeSelect: document.getElementById("resizeModeSelect"),
    optimize: document.getElementById("optimize")
  };

  // Validate required elements
  if (!elements.form || !elements.form.querySelector('input[type="file"]')) {
    console.error("未找到所需的表单元素");
    return;
  }

  const fileInput = elements.form.querySelector('input[type="file"]');
  const MAX_FILE_SIZE = 32 * 1024 * 1024; // 32MB

  // Setup quick actions
  setupQuickActions();
  
  // Setup background removal toggle
  if (elements.removeBackground) {
    elements.removeBackground.addEventListener("change", () => {
      if (elements.bgRemovalOptions) {
        elements.bgRemovalOptions.classList.toggle("hidden", !elements.removeBackground.checked);
      }
    });
  }

  // Handle form submission
  elements.form.addEventListener("submit", handleFormSubmit);

  function setupQuickActions() {
    const actions = {
      optimize: () => {
        elements.formatSelect.value = "webp";
        elements.qualitySelect.value = "medium";
        elements.widthInput.value = "";
        elements.heightInput.value = "";
      },
      resize: () => {
        elements.widthInput.focus();
      },
      adjust: () => {
        elements.formatSelect.value = "png";
        elements.qualitySelect.value = "lossless";
      },
      convert: () => elements.formatSelect.focus()
    };

    document.querySelectorAll("[data-action]").forEach(button => {
      button.addEventListener("click", (e) => {
        e.preventDefault();
        const action = button.dataset.action;
        if (actions[action]) {
          actions[action]();
          updateQuickActionStyles(button);
        }
      });
    });
  }

  async function handleFormSubmit(e) {
    e.preventDefault();

    if (!validateFile()) return;

    const formData = createFormData();
    
    try {
      await processImage(formData);
    } catch (error) {
      console.error("Error:", error);
      showError(error.message);
    }
  }

  function validateFile() {
    if (!fileInput.files || !fileInput.files.length) {
      showError("请选择一个图像文件");
      return false;
    }

    const file = fileInput.files[0];
    if (file.size > MAX_FILE_SIZE) {
      showError("文件大小超出 32MB 限制");
      return false;
    }

    return true;
  }

  function createFormData() {
    const formData = new FormData();
    formData.append("image", fileInput.files[0]);
    
    // Add processing options
    if (elements.formatSelect) {
      formData.append("format", elements.formatSelect.value || "jpeg");
    }
    if (elements.qualitySelect) {
      formData.append("quality", elements.qualitySelect.value);
    }
    if (elements.widthInput?.value) {
      formData.append("width", elements.widthInput.value);
    }
    if (elements.heightInput?.value) {
      formData.append("height", elements.heightInput.value);
    }
    if (elements.resizeModeSelect) {
      formData.append("resizeMode", elements.resizeModeSelect.value || "fit");
    }
    if (elements.removeBackground?.checked) {
      formData.append("removeBackground", "true");
    }
    if (elements.optimize?.checked) {
      formData.append("optimize", "true");
    }

    return formData;
  }

  async function processImage(formData) {
    showProgress();
    const submitButton = elements.form.querySelector('button[type="submit"]');
    if (submitButton) submitButton.disabled = true;

    try {
      const response = await fetch("/process", {
        method: "POST",
        body: formData
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error?.message || "处理失败");
      }

      await handleSuccess(response, elements.formatSelect?.value || "jpeg");
    } catch (error) {
      showError(error.message || "图像处理失败");
    } finally {
      hideProgress();
      if (submitButton) submitButton.disabled = false;
    }
  }

  function showProgress() {
    if (elements.progress) {
      elements.progress.style.display = "block";
    }
    if (elements.result) {
      elements.result.style.display = "none";
    }
  }

  function hideProgress() {
    if (elements.progress) {
      elements.progress.style.display = "none";
    }
  }

  function showError(message) {
    const errorDiv = document.createElement("div");
    errorDiv.className = "bg-red-50 border-l-4 border-red-400 p-4 mb-4";
    errorDiv.innerHTML = `
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">${message}</p>
        </div>
      </div>
    `;
    elements.form.insertBefore(errorDiv, elements.form.firstChild);
    setTimeout(() => errorDiv.remove(), 5000);
  }

  async function handleSuccess(response, format) {
    if (!elements.result) {
      console.error("未找到结果元素");
      return;
    }

    const blob = await response.blob();
    const url = URL.createObjectURL(blob);

    const resultPreview = elements.result.querySelector(".result-preview");
    const resultInfo = elements.result.querySelector(".result-info");

    if (resultPreview && resultInfo) {
      resultPreview.innerHTML = `<img src="${url}" alt="Processed image" class="max-w-full rounded-lg">`;
      resultInfo.innerHTML = `
        <div class="flex justify-between items-center">
          <p class="text-sm text-gray-600">文件大小: ${(blob.size / 1024).toFixed(2)} KB</p>
          <a href="${url}" 
             download="processed.${format}" 
             class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            下载图片
          </a>
        </div>
      `;
    }

    elements.result.style.display = "block";
  }

  function updateQuickActionStyles(button) {
    document.querySelectorAll("[data-action]").forEach((btn) => {
      btn.classList.remove(
        "bg-indigo-50",
        "border-indigo-500",
        "text-indigo-700"
      );
    });
    button.classList.add(
      "bg-indigo-50",
      "border-indigo-500",
      "text-indigo-700"
    );
  }
});
