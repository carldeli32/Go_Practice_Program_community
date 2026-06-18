<template>
  <div class="editor-wrapper border border-border rounded-lg overflow-hidden bg-input"
    :class="{ 'ring-1 ring-primary': editor?.isFocused }">
    <!-- Toolbar -->
    <div class="flex items-center gap-0.5 px-2 py-1.5 border-b border-border bg-card flex-wrap">
      <button v-for="btn in toolbarBtns" :key="btn.label"
        @click="btn.action"
        :title="btn.label"
        :class="['w-7 h-7 flex items-center justify-center rounded text-xs text-muted-fg hover:text-fg hover:bg-secondary transition-colors cursor-pointer bg-transparent border-0', btn.isActive?.() ? 'text-primary bg-primary/10' : '']"
      >
        <component :is="btn.icon" :size="14" />
      </button>

      <span class="w-px h-5 bg-border mx-1"></span>

      <!-- Image upload -->
      <button @click="$refs.fileInput.click()" title="插入图片"
        :disabled="uploading"
        class="w-7 h-7 flex items-center justify-center rounded text-xs text-muted-fg hover:text-primary hover:bg-secondary transition-colors cursor-pointer bg-transparent border-0 disabled:opacity-40"
      >
        <Loader v-if="uploading" :size="14" class="animate-spin" />
        <ImageIcon v-else :size="14" />
      </button>
      <input ref="fileInput" type="file" accept="image/*" hidden @change="handleImageUpload" />
    </div>

    <!-- Editor -->
    <EditorContent :editor="editor" class="prose prose-sm max-w-none px-3 py-2 min-h-[120px] text-sm text-fg" />
  </div>
</template>

<script setup>
import { ref, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import { Bold, Italic, Heading2, List, ListOrdered, Quote, Undo2, Redo2, ImageIcon, Loader } from 'lucide-vue-next'
import api from '../api'
import { useToast } from '../composables/toast'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '写下你的想法...' },
  postId: { type: Number, default: 0 },
})
const emit = defineEmits(['update:modelValue'])

const toast = useToast()
const uploading = ref(false)

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      heading: { levels: [2, 3] },
    }),
    Image.configure({ inline: false, allowBase64: false }),
  ],
  editorProps: {
    attributes: {
      placeholder: props.placeholder,
    },
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
})

// Watch external modelValue changes
watch(() => props.modelValue, (val) => {
  if (editor.value && val !== editor.value.getHTML()) {
    editor.value.commands.setContent(val, false)
  }
})

const toolbarBtns = [
  { label: '加粗', icon: Bold, action: () => editor.value?.chain().focus().toggleBold().run(), isActive: () => editor.value?.isActive('bold') },
  { label: '斜体', icon: Italic, action: () => editor.value?.chain().focus().toggleItalic().run(), isActive: () => editor.value?.isActive('italic') },
  { label: '标题', icon: Heading2, action: () => editor.value?.chain().focus().toggleHeading({ level: 2 }).run(), isActive: () => editor.value?.isActive('heading', { level: 2 }) },
  { label: '无序列表', icon: List, action: () => editor.value?.chain().focus().toggleBulletList().run(), isActive: () => editor.value?.isActive('bulletList') },
  { label: '有序列表', icon: ListOrdered, action: () => editor.value?.chain().focus().toggleOrderedList().run(), isActive: () => editor.value?.isActive('orderedList') },
  { label: '引用', icon: Quote, action: () => editor.value?.chain().focus().toggleBlockquote().run(), isActive: () => editor.value?.isActive('blockquote') },
  { label: '撤销', icon: Undo2, action: () => editor.value?.chain().focus().undo().run() },
  { label: '重做', icon: Redo2, action: () => editor.value?.chain().focus().redo().run() },
]

async function handleImageUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return

  uploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    form.append('post_id', props.postId)
    const res = await api.post('/upload/image', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    const url = res.data.data.url
    editor.value?.chain().focus().setImage({ src: url }).run()
  } catch (err) {
    toast.error('上传失败，图片不能超过 5MB')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style>
/* Tiptap editor styles */
.editor-wrapper .ProseMirror {
  outline: none;
  min-height: 120px;
  line-height: 1.6;
}
.editor-wrapper .ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  color: var(--color-muted-fg);
  pointer-events: none;
  float: left;
  height: 0;
}
.editor-wrapper .ProseMirror img {
  max-width: 100%;
  border-radius: 8px;
  margin: 8px 0;
}
.editor-wrapper .ProseMirror blockquote {
  border-left: 3px solid var(--color-primary);
  padding-left: 12px;
  color: var(--color-muted-fg);
  margin: 8px 0;
}
.editor-wrapper .ProseMirror h2 {
  font-size: 18px;
  font-weight: 700;
  font-family: var(--font-heading);
  margin: 12px 0 4px;
}
.editor-wrapper .ProseMirror ul,
.editor-wrapper .ProseMirror ol {
  padding-left: 20px;
  margin: 4px 0;
}
</style>
