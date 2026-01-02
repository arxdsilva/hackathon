require("expose-loader?exposes=$,jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");
require("jquery-ujs/src/rails.js");
require("easymde/dist/easymde.min.css");

const EasyMDE = require("easymde");
const { marked } = require("marked");
const createDOMPurify = require("dompurify");

$(() => {
	const DOMPurify = createDOMPurify(window);

	const attachMarkdownEditors = () => {
		document.querySelectorAll("[data-markdown-editor]").forEach((textarea) => {
			const required = textarea.dataset.markdownRequired === "true";
			const requiredMessage = textarea.dataset.markdownRequiredMessage || "Please fill out this field.";

			const editor = new EasyMDE({
				element: textarea,
				spellChecker: false,
				status: false,
				toolbar: [
					"bold",
					"italic",
					"heading",
					"|",
					"quote",
					"unordered-list",
					"ordered-list",
					"|",
					"link",
					"image",
					"code",
					"table",
					"horizontal-rule",
					"|",
					"preview",
					"side-by-side",
					"fullscreen",
				],
			});

			if (textarea.form) {
				textarea.form.addEventListener("submit", (event) => {
					const value = editor.value().trim();
					textarea.value = value;

					if (required && value === "") {
						event.preventDefault();
						event.stopImmediatePropagation();
						window.alert(requiredMessage);
						editor.codemirror.focus();
					}
				});
			}
		});
	};

	const renderMarkdown = () => {
		document.querySelectorAll("[data-markdown-source]").forEach((node) => {
			const source = (node.textContent || "").trim();
			const html = DOMPurify.sanitize(marked.parse(source));
			node.innerHTML = html;
		});
	};

	attachMarkdownEditors();
	renderMarkdown();
});
