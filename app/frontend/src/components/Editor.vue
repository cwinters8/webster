<script setup lang="ts">
import { onMounted, ref } from "vue"
import "tinymce"
import Editor from "@tinymce/tinymce-vue"

// required components
import "tinymce/icons/default/icons.min.js"
import "tinymce/themes/silver/theme.min.js"
import "tinymce/models/dom/model.min.js"

// UI styling
import "tinymce/skins/ui/oxide/skin.js"

// plugins
import "tinymce/plugins/advlist"
import "tinymce/plugins/code"
import "tinymce/plugins/emoticons"
import "tinymce/plugins/emoticons/js/emojis"
import "tinymce/plugins/link"
import "tinymce/plugins/lists"
import "tinymce/plugins/table"

// content styling
import "tinymce/skins/content/default/content.js"

const props = defineProps<{
  id?: string
}>()

const html = ref("")

onMounted(() => {
  if (props.id && props.id.length > 0) {
    // TODO: call Go function to retrieve content
  }
})
</script>

<template>
  <div class="container">
    <Editor
      :init="{ height: '100vh' }"
      initial-value="<h1>Hello, Webster</h1>"
      v-model="html"
    />
  </div>
</template>

<style scoped></style>
