package langserver

import (
	"fmt"
	"strings"
)

const wildcard = "="

// DocumentURI is the uri of a document(a file or a folder)
type DocumentURI string

// InitializeParams is
type InitializeParams struct {
	// ProcessID is the process Id of the parent process that started the server.
	// Is null if the process has not been started by another process.
	// If the parent process is not alive then the server should exit (see exit notification) its process.
	ProcessID int `json:"processId,omitempty"`
	// RootURI is the rootUri of the workspace. Is null if no folder is open. If both `rootPath` and `rootUri` are set `rootUri` wins.
	// @deprecated in favour of `workspaceFolders`
	RootURI DocumentURI `json:"rootUri,omitempty"`
	// InitializationOptions is user provided initialization options
	InitializationOptions InitializeOptions `json:"initializationOptions,omitempty"`
	// Capabilities are the capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities,omitempty"`
	// Trace is the initial trace setting. If omitted trace is disabled ('off').
	Trace string `json:"trace,omitempty"`
}

// InitializeOptions is
type InitializeOptions struct {
	DocumentFormatting bool `json:"documentFormatting"`
	Hover              bool `json:"hover"`
	DocumentSymbol     bool `json:"documentSymbol"`
	CodeAction         bool `json:"codeAction"`
	Completion         bool `json:"completion"`
}

// ClientCapabilities is the workspace specific client capabilities
type ClientCapabilities struct {
}

// InitializeResult is
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

// MessageType is
type MessageType int

// LogError is
const (
	LogError   MessageType = 1
	LogWarning MessageType = 2
	LogInfo    MessageType = 3
	LogLog     MessageType = 4
)

// TextDocumentSyncKind defines how the host (editor) should sync document changes to the language server.
type TextDocumentSyncKind int

// TDSKNone is
const (
	// TDSKNone means documents should not be synced at all.
	TDSKNone TextDocumentSyncKind = 0
	// TDSKFull means documents are synced by always sending the full content of the document.
	TDSKFull TextDocumentSyncKind = 1
	// TDSKIncremental means documents are synced by sending the full content on open. After that only incremental updates to the document
	// are send.
	TDSKIncremental TextDocumentSyncKind = 2
)

// CompletionProvider is
type CompletionProvider struct {
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
	TriggerCharacters []string `json:"triggerCharacters"`
}

// WorkspaceFoldersServerCapabilities is
type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported"`
	ChangeNotifications bool `json:"changeNotifications"`
}

// ServerCapabilitiesWorkspace is
type ServerCapabilitiesWorkspace struct {
	WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders"`
}

// ServerCapabilities is the capabilities the language server provides
// see: https://microsoft.github.io/language-server-protocol/specifications/specification-3-17/#serverCapabilities
type ServerCapabilities struct {
	// TextDocumentSync defines how text documents are synced.
	// Is either a detailed structure defining each notification or for backwards compatibility the TextDocumentSyncKind number.
	// If omitted it defaults to `TextDocumentSyncKind.None`.
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync,omitempty"`
	// DocumentSymbolProvider is true if the server provides document symbol support.
	DocumentSymbolProvider bool `json:"documentSymbolProvider,omitempty"`
	// CompletionProvider defines how the server provides completion support.
	CompletionProvider *CompletionProvider `json:"completionProvider,omitempty"`
	// DefinitionProvider is true if the server provides goto definition support.
	DefinitionProvider bool `json:"definitionProvider,omitempty"`
	// ReferencesProvider is true if the server provides find references support.
	ReferencesProvider bool `json:"referencesProvider,omitempty"`
	// DocumentFormattingProvider is true if the server provides document formatting.
	DocumentFormattingProvider bool `json:"documentFormattingProvider,omitempty"`
	// HoverProvider is true if the server provides hover support.
	HoverProvider bool `json:"hoverProvider,omitempty"`
	// CodeActionProvider is true if the server provides code actions.
	// The `CodeActionOptions` return type is only valid if the client signals code action literal support via the property
	// `textDocument.codeAction.codeActionLiteralSupport`.
	CodeActionProvider bool `json:"codeActionProvider,omitempty"`
	// Workspace specific server capabilities
	Workspace *ServerCapabilitiesWorkspace `json:"workspace,omitempty"`
}

// TextDocumentItem is
type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

// VersionedTextDocumentIdentifier is
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentIdentifier is the identifier of a text document. Text documents are identified using a URI.
type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

// DidOpenTextDocumentParams is
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidCloseTextDocumentParams is
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// TextDocumentContentChangeEvent is
type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range"`
	RangeLength int    `json:"rangeLength"`
	Text        string `json:"text"`
}

// DidChangeTextDocumentParams is
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidSaveTextDocumentParams is
type DidSaveTextDocumentParams struct {
	Text         *string                `json:"text"`
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// TextDocumentPositionParams is
type TextDocumentPositionParams struct {
	// TextDocument is the text document
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Position is the position inside the text document.
	Position Position `json:"position"`
}

// CompletionParams is the params of a textDocument/completion request
// see https://microsoft.github.io/language-server-protocol/specifications/specification-3-17/#completionParams
type CompletionParams struct {
	TextDocumentPositionParams
	// CompletionContext is the completion context. This is only available if the client specifies to send this using the client capability
	// `completion.contextSupport === true`
	CompletionContext CompletionContext `json:"contentChanges"`
}

// CompletionContext contains additional information about the context in which a completion request is triggered.
type CompletionContext struct {
	// TriggerKind is how the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind"`
	// TriggerCharacter is the trigger character (a single character) that has trigger code complete.
	// Is undefined if `triggerKind !== CompletionTriggerKind.TriggerCharacter`
	TriggerCharacter *string `json:"triggerCharacter"`
}

type CompletionTriggerKind int

const (
	// Invoked means the completion was triggered by typing an identifier (24x7 code complete), manual invocation (e.g Ctrl+Space) or via API.
	Invoked CompletionTriggerKind = 1
	// TriggerCharacter means the completion was triggered by a trigger character specified by the `triggerCharacters` properties of the
	// `CompletionRegistrationOptions`.
	TriggerCharacter CompletionTriggerKind = 2
	// TriggerForIncompleteCompletions means the completion was re-triggered as the current completion list is incomplete.
	TriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// HoverParams is
type HoverParams struct {
	TextDocumentPositionParams
}

// Location is
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// Range is
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position is between two characters like an ‘insert’ cursor in a editor. Position in a text document is expressed as zero-based line and
// zero-based character offset. Special values like for example -1 to denote the end of a line are not supported.
type Position struct {
	// Line means the line position in a document (zero-based).
	Line int `json:"line"`
	// Character means the character offset on a line in a document (zero-based).
	// Assuming that the line is represented as a string, the `character` value represents the gap between the `character` and `character + 1`.
	// If the character value is greater than the line length it defaults back to the line length.
	Character int `json:"character"`
}

// DiagnosticRelatedInformation is
type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// Diagnostic is
type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           int                            `json:"severity,omitempty"`
	Code               *string                        `json:"code,omitempty"`
	Source             *string                        `json:"source,omitempty"`
	Message            string                         `json:"message"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

// PublishDiagnosticsParams is
type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Version     int          `json:"version"`
}

// FormattingOptions is
type FormattingOptions map[string]interface{}

// DocumentFormattingParams is
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// TextEdit is
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// DocumentSymbolParams is
type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentSymbol is
type DocumentSymbol struct {
	Name           string           `json:"name"`
	Kind           int64            `json:"kind"`
	Deprecated     bool             `json:"deprecated"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
	Detail         string           `json:"detail,omitempty"`
}

// SymbolInformation is
type SymbolInformation struct {
	Name          string   `json:"name"`
	Kind          int64    `json:"kind"`
	Deprecated    bool     `json:"deprecated"`
	Location      Location `json:"location"`
	ContainerName *string  `json:"containerName"`
}

// CompletionItemKind is
type CompletionItemKind int

// TextCompletion is
const (
	TextCompletion          CompletionItemKind = 1
	MethodCompletion        CompletionItemKind = 2
	FunctionCompletion      CompletionItemKind = 3
	ConstructorCompletion   CompletionItemKind = 4
	FieldCompletion         CompletionItemKind = 5
	VariableCompletion      CompletionItemKind = 6
	ClassCompletion         CompletionItemKind = 7
	InterfaceCompletion     CompletionItemKind = 8
	ModuleCompletion        CompletionItemKind = 9
	PropertyCompletion      CompletionItemKind = 10
	UnitCompletion          CompletionItemKind = 11
	ValueCompletion         CompletionItemKind = 12
	EnumCompletion          CompletionItemKind = 13
	KeywordCompletion       CompletionItemKind = 14
	SnippetCompletion       CompletionItemKind = 15
	ColorCompletion         CompletionItemKind = 16
	FileCompletion          CompletionItemKind = 17
	ReferenceCompletion     CompletionItemKind = 18
	FolderCompletion        CompletionItemKind = 19
	EnumMemberCompletion    CompletionItemKind = 20
	ConstantCompletion      CompletionItemKind = 21
	StructCompletion        CompletionItemKind = 22
	EventCompletion         CompletionItemKind = 23
	OperatorCompletion      CompletionItemKind = 24
	TypeParameterCompletion CompletionItemKind = 25
)

// CompletionItemTag is
type CompletionItemTag int

// InsertTextFormat is
type InsertTextFormat int

// PlainTextTextFormat is
const (
	PlainTextTextFormat InsertTextFormat = 1
	SnippetTextFormat   InsertTextFormat = 2
)

// Command is
type Command struct {
	Title     string        `json:"title" yaml:"title"`
	Command   string        `json:"command" yaml:"command"`
	Arguments []interface{} `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	OS        string        `json:"-" yaml:"os,omitempty"`
}

// WorkspaceEdit is
type WorkspaceEdit struct {
	Changes         interface{} `json:"changes"`         // { [uri: DocumentUri]: TextEdit[]; };
	DocumentChanges interface{} `json:"documentChanges"` // (TextDocumentEdit[] | (TextDocumentEdit | CreateFile | RenameFile | DeleteFile)[]);
}

// CodeAction is
type CodeAction struct {
	Title       string         `json:"title"`
	Diagnostics []Diagnostic   `json:"diagnostics"`
	IsPreferred bool           `json:"isPreferred"` // TODO
	Edit        *WorkspaceEdit `json:"edit"`
	Command     *Command       `json:"command"`
}

// CompletionItem is
// see: https://microsoft.github.io/language-server-protocol/specifications/specification-3-17/#completionItem
type CompletionItem struct {
	// Label is the label of this completion item. By default also the text that is inserted when selecting this completion.
	Label string `json:"label"`
	// Kind is the kind of this completion item. Based of the kind an icon is chosen by the editor.
	// The standardized set of available values is defined in `CompletionItemKind`.
	Kind CompletionItemKind `json:"kind,omitempty"`
	// Tags is the Tags for this completion item.
	// @since 3.15.0
	Tags []CompletionItemTag `json:"tags,omitempty"`
	// Detail is a human-readable string with additional information about this item, like type or symbol information.
	Detail string `json:"detail,omitempty"`
	// Documentation is a human-readable string that represents a doc-comment.
	Documentation string `json:"documentation,omitempty"` // string | MarkupContent
	// Deprecated indicates if this item is deprecated.
	Deprecated bool `json:"deprecated,omitempty"`
	// Preselect defines whether to select this item when showing.
	// *Note* that only one completion item can be selected and that the tool / client decides which item that is.
	// The rule is that the *first* item of those that match best is selected.
	Preselect bool `json:"preselect,omitempty"`
	// SortText is a string that should be used when comparing this item with other items.
	// When `falsy` the label is used as the sort text for this item.
	SortText string `json:"sortText,omitempty"`
	// FilterText is a string that should be used when filtering a set of completion items.
	// When `falsy` the label is used as the filter text for this item.
	FilterText string `json:"filterText,omitempty"`
	// InsertText is a string that should be inserted into a document when selecting this completion.
	// When `falsy` the label is used as the insert text for this item.
	//
	// The `insertText` is subject to interpretation by the client side. Some tools might not take the string literally.
	// For example VS Code when code complete is requested in this example `con<cursor position>` and a completion item with an `insertText`
	// of `console` is provided it will only insert `sole`. Therefore it is recommended to use `textEdit` instead since it avoids additional
	// client side interpretation.
	InsertText string `json:"insertText,omitempty"`
	// InsertTextFormat is the format of the insert text.
	// The format applies to both the `insertText` property and the `newText` property of a provided `textEdit`.
	// If omitted defaults to `InsertTextFormat.PlainText`.
	InsertTextFormat InsertTextFormat `json:"insertTextFormat,omitempty"`
	// InsertTextMode defines how whitespace and indentation is handled during completion item insertion.
	// If not provided the client's default value is used.
	// @since 3.16.0
	// insertTextMode InsertTextMode `json:"insertTextMode,omitempty`

	// TextEdit is the edit which is applied to a document when selecting this completion.
	// When an edit is provided the value of `insertText` is ignored.
	//
	// *Note:* The range of the edit must be a single line range and it must contain the position at which completion has been requested.
	// Most editors support two different operations when accepting a completion item.
	// One is to insert a completion text and the other is to replace an existing text with a completion text.
	// Since this can usually not be predetermined by a server it can report both ranges.
	// Clients need to signal support for `InsertReplaceEdit`s via the `textDocument.completion.insertReplaceSupport` client capability property.
	//
	// *Note 1:* The text edit's range as well as both ranges from an insert replace edit must be a [single line] and they must contain the
	// position at which completion has been requested.
	// *Note 2:* If an `InsertReplaceEdit` is returned the edit's insert range must be a prefix of the edit's replace range, that means it
	// must be contained and starting at the same position.
	//
	// @since 3.16.0 additional type `InsertReplaceEdit`
	TextEdit *TextEdit `json:"textEdit,omitempty"`
	// AdditionalTextEdits is an optional array of additional text edits that are applied when selecting this completion.
	// Edits must not overlap (including the same insert position) with the main edit nor with themselves.
	//
	// Additional text edits should be used to change text unrelated to the current cursor position (for example adding an import statement
	// at the top of the file if the completion item will insert an unqualified type).
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
	// CommitCharacters is an optional set of characters that when pressed while this completion is active will accept it first and then
	// type that character.
	// *Note* that all commit characters should have `length=1` and that superfluous characters will be ignored.
	CommitCharacters []string `json:"commitCharacters,omitempty"`
	// Command is an optional command that is executed *after* inserting this completion.
	// *Note* that additional modifications to the current document should be described with the additionalTextEdits-property.
	Command *Command `json:"command,omitempty"`
	// Data is a data entry field that is preserved on a completion item between a completion and a completion resolve request.
	Data interface{} `json:"data,omitempty"`
}

// Hover is
type Hover struct {
	Contents interface{} `json:"contents"`
	Range    *Range      `json:"range"`
}

// MarkedString is
type MarkedString struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

// MarkupKind is
type MarkupKind string

// PlainText is
const (
	PlainText MarkupKind = "plaintext"
	Markdown  MarkupKind = "markdown"
)

// MarkupContent is
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

// WorkDoneProgressParams is
type WorkDoneProgressParams struct {
	WorkDoneToken interface{} `json:"workDoneToken"`
}

// ExecuteCommandParams is
type ExecuteCommandParams struct {
	WorkDoneProgressParams

	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

// CodeActionKind is
type CodeActionKind string

// Empty is
const (
	Empty                 CodeActionKind = ""
	QuickFix              CodeActionKind = "quickfix"
	Refactor              CodeActionKind = "refactor"
	RefactorExtract       CodeActionKind = "refactor.extract"
	RefactorInline        CodeActionKind = "refactor.inline"
	RefactorRewrite       CodeActionKind = "refactor.rewrite"
	Source                CodeActionKind = "source"
	SourceOrganizeImports CodeActionKind = "source.organizeImports"
)

// CodeActionContext is
type CodeActionContext struct {
	Diagnostics []Diagnostic     `json:"diagnostics"`
	Only        []CodeActionKind `json:"only,omitempty"`
}

// PartialResultParams is
type PartialResultParams struct {
	PartialResultToken interface{} `json:"partialResultToken"`
}

// CodeActionParams is
type CodeActionParams struct {
	WorkDoneProgressParams
	PartialResultParams

	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

// DidChangeConfigurationParams is
type DidChangeConfigurationParams struct {
	Settings Config `json:"settings"`
}

// NotificationMessage is
type NotificationMessage struct {
	Method string      `json:"message"`
	Params interface{} `json:"params"`
}

// DocumentDefinitionParams is
type DocumentDefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// ReferenceParams is
type ReferenceParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// ShowMessageParams is
type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// LogMessageParams is
type LogMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// DidChangeWorkspaceFoldersParams is
type DidChangeWorkspaceFoldersParams struct {
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// WorkspaceFoldersChangeEvent is
type WorkspaceFoldersChangeEvent struct {
	Added   []WorkspaceFolder `json:"added,omitempty"`
	Removed []WorkspaceFolder `json:"removed,omitempty"`
}

// WorkspaceFolder is
type WorkspaceFolder struct {
	URI  DocumentURI `json:"uri"`
	Name string      `json:"name"`
}

// PathToURI convert fileName to fileURI.
// the provided fileName should be a valid absolute file path
func PathToURI(fileName string) string {
	return fmt.Sprintf("file://%s", fileName)
}

const fileURIPrefix = "file://"

// URIToPath convert fileURI to fileName.
// the provided fileURI should be a valid
func URIToPath(uri DocumentURI) (string, error) {
	if !strings.HasPrefix(string(uri), fileURIPrefix) {
		return "", nil
	}
	return strings.Replace(string(uri), fileURIPrefix, "", 1), nil
}
