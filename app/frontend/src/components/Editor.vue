<script setup>
import Editor from "@tinymce/tinymce-vue"
</script>

<template>
  <div class="container">
    <Editor :init="{ height: '100vh' }"><h1>Hello, Webster!</h1></Editor>
  </div>
</template>

<style scoped>
/* .container {
  border: 5px black;
  border-radius: 10px;
} */
</style>
