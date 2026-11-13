<script setup lang="ts">
import { MIT_CREDENTIAL_URL } from '~/utils/credentials'
</script>

<template>
  <!-- Accredible credential (details from Nick, 2026-11-11). -->
  <div class="credential">
    <div class="credential__seal" aria-hidden="true">MIT<br />PE</div>
    <div>
      <strong>MIT Professional Education</strong>
      <p>Applied Agentic AI for Organizational Transformation · verified on Accredible</p>
      <a :href="MIT_CREDENTIAL_URL" target="_blank" rel="noopener">View credential</a>
    </div>
  </div>
</template>
