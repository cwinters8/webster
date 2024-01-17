<script setup>
import { useEditor, EditorContent } from "@tiptap/vue-3"
import StarterKit from "@tiptap/starter-kit"

const editor = useEditor({
  content: "<h1>Hello, Webster!</h1>",
  extensions: [StarterKit],
})

console.log({ html: editor.value.getHTML() })
</script>

<template>
  <editor-content :editor="editor" />
</template>
