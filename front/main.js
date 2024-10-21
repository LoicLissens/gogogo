import htmx from "htmx.org";
import { BooleanCheckbox } from "./custom_elements.js";

customElements.define("boolean-checkbox", BooleanCheckbox, { extends: "input" });