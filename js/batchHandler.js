document.addEventListener("DOMContentLoaded", function () {
    const elements = {
        batchUploadArea: document.getElementById("batchUploadArea"),
        batchImageInput: document.getElementById("batchImageInput"),
        batchFileList: document.getElementById("batchFileList"),
        batchProcessBtn: document.getElementById("batchProcessBtn"),
        batchProgress: document.getElementById("batchProgress"),
        batchProgressBar: document.getElementById("batchProgressBar"),
        batchProgressCount: document.getElementById("batchProgressCount")
    };

    const state = {
        files: [],
        processing: false
    };

    function initializeBatchUpload() {
        if (!elements.batchUploadArea || !elements.batchImageInput) {
            console.error("未找到所需的批量上传元素");
            return;
        }

        setupBatchEventListeners();
    }

    function setupBatchEventListeners() {
        // File input change handler
        elements.batchImageInput.addEventListener("change", handleFileSelect);

        // Drag and drop handlers
        elements.batchUploadArea.addEventListener("dragover", (e) => {
            e.preventDefault();
            elements.batchUploadArea.classList.add("border-indigo-500");
        });

        elements.batchUploadArea.addEventListener("dragleave", (e) => {
            e.preventDefault();
            elements.batchUploadArea.classList.remove("border-indigo-500");
        });

        elements.batchUploadArea.addEventListener("drop", (e) => {
            e.preventDefault();
            elements.batchUploadArea.classList.remove("border-indigo-500");
            
            if (e.dataTransfer.files.length) {
                handleFileSelect({ target: { files: e.dataTransfer.files } });
            }
        });

        // Process button handler
        elements.batchProcessBtn.addEventListener("click", processBatch);
    }

    function handleFileSelect(e) {
        const files = Array.from(e.target.files).filter(file => {
            const isImage = file.type.startsWith("image/");
            const isHeic = file.name.toLowerCase().endsWith('.heic') || 
                          file.name.toLowerCase().endsWith('.heif');
            return isImage || isHeic;
        });
        
        if (files.length === 0) {
            alert("请选择有效的图像文件");
            return;
        }
    
        state.files = files;
        updateFileList();
        elements.batchProcessBtn.disabled = false;
    }

    function updateFileList() {
        elements.batchFileList.innerHTML = state.files.map((file, index) => `
            <div class="flex items-center justify-between py-2">
                <div class="flex items-center">
                    <span class="text-sm font-medium transition-colors" 
                          :class="{ 'text-darkTextPrimary': darkMode, 'text-gray-900': !darkMode }">
                        ${file.name}
                    </span>
                    <span class="ml-2 text-sm transition-colors" 
                          :class="{ 'text-darkTextSecondary': darkMode, 'text-gray-500': !darkMode }">
                        (${(file.size / (1024 * 1024)).toFixed(2)} MB)
                    </span>
                </div>
                <button onclick="removeFile(${index})" 
                        class="transition-colors hover:text-red-700"
                        :class="{ 'text-red-400': darkMode, 'text-red-500': !darkMode }">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
        `).join("");
        
        elements.batchFileList.classList.remove("hidden");
    }

    window.removeFile = function(index) {
        state.files.splice(index, 1);
        updateFileList();
        elements.batchProcessBtn.disabled = state.files.length === 0;
    };

    async function processBatch() {
        if (state.processing || state.files.length === 0) return;

        state.processing = true;
        elements.batchProcessBtn.disabled = true;
        elements.batchProgress.classList.remove('hidden');

        const options = getProcessingOptions();

        if (options.mode === 'merge-pdf') {
            await processMergePDF(options);
        } else {
            await processIndividualFiles(options);
        }

        state.processing = false;
        elements.batchProcessBtn.disabled = false;
        resetBatchUpload();
    }

    function getProcessingOptions() {
        const processingMode = document.querySelector('input[name="processing-mode"]:checked').value;
        
        if (processingMode === 'merge-pdf') {
            return {
                mode: 'merge-pdf',
                pageSize: document.getElementById('pdfPageSize').value,
                orientation: document.getElementById('pdfOrientation').value,
                imagesPerPage: document.getElementById('imagesPerPage').value
            };
        }
        
        return {
            mode: 'individual',
            format: document.getElementById('batchFormatSelect')?.value || 'jpeg',
            // quality: document.getElementById('batchQualitySelect')?.value || 'medium',
            width: document.getElementById('batchWidthInput')?.value || '',
            height: document.getElementById('batchHeightInput')?.value || '',
            // optimize: document.getElementById('batchOptimize')?.checked || false
        };
    }

    async function processIndividualFiles(options) {
        const total = state.files.length;
        let processed = 0;
        let errors = [];

        try {
            for (const file of state.files) {
                const formData = new FormData();
                formData.append("image", file);
                
                // Add processing options
                for (const [key, value] of Object.entries(options)) {
                    formData.append(key, value);
                }

                try {
                    await processFile(formData);
                    processed++;
                } catch (error) {
                    errors.push(`处理失败 ${file.name}: ${error.message}`);
                }
                
                updateProgress(processed, total);
            }

            if (errors.length > 0) {
                alert(`批量处理完成，但存在 ${errors.length} 个错误:\n${errors.join('\n')}`);
            } else {
                alert("批量处理成功完成！");
            }
        } catch (error) {
            console.error("批量处理错误:", error);
            alert("批量处理期间发生错误");
        }
    }

    async function processMergePDF(options) {
        const formData = new FormData();
        state.files.forEach((file, index) => {
            formData.append(`images`, file);
        });
        
        // Add PDF options
        formData.append('mode', 'merge-pdf');
        formData.append('pageSize', options.pageSize);
        formData.append('orientation', options.orientation);
        formData.append('imagesPerPage', options.imagesPerPage);

        try {
            const response = await fetch('/process/merge-pdf', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('PDF 创建失败');
            }

            const blob = await response.blob();
            downloadFile(blob, `merged_${Date.now()}.pdf`);
        } catch (error) {
            console.error('创建 PDF 时出错:', error);
            alert('创建 PDF 失败: ' + error.message);
        }
    }

    async function processFile(formData) {
        try {
            const response = await fetch("/process", {
                method: "POST",
                body: formData
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || "处理失败");
            }

            const blob = await response.blob();
            const originalName = formData.get("image").name;
            const format = formData.get("format") || "jpeg";
            
            // Get the original file extension
            const originalExt = originalName.split('.').pop();
            
            // Create the new filename
            const baseName = originalName.replace(`.${originalExt}`, '');
            let filename = `processed_${baseName}.${format}`;
            
            // Set the correct content type based on the format
            const contentType = getContentType(format);
            const processedBlob = new Blob([blob], { type: contentType });
            
            downloadFile(processedBlob, filename);
            return true;
        } catch (error) {
            console.error(`处理文件时出错: ${error.message}`);
            throw error;
        }
    }

    function getContentType(format) {
        const contentTypes = {
            'jpeg': 'image/jpeg',
            'jpg': 'image/jpeg',
            'png': 'image/png',
            'webp': 'image/webp',
            'gif': 'image/gif',
            'bmp': 'image/bmp',
            'heic': 'image/heic',
            'heif': 'image/heif',
            'ico': 'image/x-icon',
            'pdf': 'application/pdf'
        };
        return contentTypes[format] || 'application/octet-stream';
    }

    function downloadFile(blob, filename) {
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    function updateProgress(processed, total) {
        const percentage = (processed / total) * 100;
        elements.batchProgressBar.style.width = `${percentage}%`;
        elements.batchProgressCount.textContent = `${processed}/${total} files`;
    }

    function resetBatchUpload() {
        state.files = [];
        elements.batchImageInput.value = "";
        elements.batchFileList.classList.add("hidden");
        elements.batchProgress.classList.add("hidden");
        elements.batchProgressBar.style.width = "0%";
        elements.batchProgressCount.textContent = "0/0 files";
        elements.batchProcessBtn.disabled = true;
    }

    initializeBatchUpload();
});