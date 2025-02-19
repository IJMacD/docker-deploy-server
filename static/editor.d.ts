/**
 * @typedef Token
 * @property {string} type
 * @property {number} start
 * @property {number} length
 */
/**
 * @param {string | Element} element
 * @param {{ tokenizer: (text: string) => Token[], plugins?: ((textarea: HTMLTextAreaElement, editor: HTMLDivElement) => () => void)[] }} options
 */
export function editor(element: string | Element, { tokenizer, plugins }: {
    tokenizer: (text: string) => Token[];
    plugins?: ((textarea: HTMLTextAreaElement, editor: HTMLDivElement) => () => void)[];
}): {
    destroy(): void;
};
/**
 * @param {string | Element} element
 * @param {((textarea: HTMLTextAreaElement, editor: HTMLDivElement) => () => void)[] } features
 */
export function baseEditor(element: string | Element, features: ((textarea: HTMLTextAreaElement, editor: HTMLDivElement) => () => void)[]): {
    destroy(): void;
};
/**
 * @param {HTMLTextAreaElement} textarea
 * @param {HTMLDivElement} editor
 */
export function resizePlugin(textarea: HTMLTextAreaElement, editor: HTMLDivElement): () => void;
/**
 * Uses a polling implementation
 * @param {HTMLTextAreaElement} textarea
 * @param {HTMLDivElement} editor
 */
export function positionPlugin(textarea: HTMLTextAreaElement, editor: HTMLDivElement): () => void;
/**
 * @param {HTMLTextAreaElement} textarea
 * @param {HTMLDivElement} editor
 */
export function scrollPlugin(textarea: HTMLTextAreaElement, editor: HTMLDivElement): () => void;
/**
 * Handles inserting tabs (as four spaces) and removing leading tabs (with
 * shift-tab).
 * @param {HTMLTextAreaElement} textarea
 */
export function tabPlugin(textarea: HTMLTextAreaElement): () => void;
/**
 * Maintains indentation on enter key
 * @param {HTMLTextAreaElement} textarea
 */
export function enterPlugin(textarea: HTMLTextAreaElement): () => void;
/**
 * Add single or double quotes, or parentheses or brackets around selected text.
 * @param {HTMLTextAreaElement} textarea
 */
export function quotesPlugin(textarea: HTMLTextAreaElement): () => void;
/**
 * Editor inherits font properties from the textarea.
 * @param {HTMLTextAreaElement} textarea
 * @param {HTMLDivElement} editor
 */
export function fontPlugin(textarea: HTMLTextAreaElement, editor: HTMLDivElement): void;
export function tokenizerPlugin(tokenizer: (text: string) => Token[]): (textarea: HTMLTextAreaElement, editor: HTMLDivElement) => () => void;
export type Token = {
    type: string;
    start: number;
    length: number;
};
//# sourceMappingURL=editor.d.ts.map