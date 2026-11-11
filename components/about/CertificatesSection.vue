<script setup lang="ts">
const certificates = [
  {
    name: 'Professional Cloud Architect',
    imgSrc: '/certs/cert_ca.webp',
    credlyUrl: 'https://www.credly.com/users/nikolay-tonev.a60be769/badges',
  },
  {
    name: 'Professional Data Engineer',
    imgSrc: '/certs/cert_de.webp',
    credlyUrl: 'https://www.credly.com/users/nikolay-tonev.a60be769/badges',
  },
  {
    name: 'Professional Developer',
    imgSrc: '/certs/cert_dev.webp',
    credlyUrl: 'https://www.credly.com/users/nikolay-tonev.a60be769/badges',
  },
  {
    name: 'Extend Google Agentspace assistant capabilities',
    imgSrc: '/certs/badge.webp',
    credlyUrl: 'https://www.credly.com/users/nikolay-tonev.a60be769/badges',
  },
  {
    name: 'Develop Advanced Enterprise Search and Conversation Applications Skill Badge',
    imgSrc: '/certs/badge2.webp',
    credlyUrl: 'https://www.credly.com/users/nikolay-tonev.a60be769/badges',
  },
  // You can add more certificates here by duplicating the object above
]
</script>

<template>
  <section class="py-20 md:py-28 bg-light-gray">
    <div class="container mx-auto px-6">
      <div class="text-center max-w-3xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-bold text-dark-slate">Certifications & Credentials</h2>
        <p class="text-lg text-muted mt-4">
          A testament to my commitment to excellence and continuous learning in the field of cloud
          and data architecture.
        </p>
      </div>

      <div class="mt-16 flex flex-wrap justify-center items-center gap-8 md:gap-16">
        <a
          v-for="cert in certificates"
          :key="cert.name"
          :href="cert.credlyUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="group"
        >
          <NuxtImg
            :src="cert.imgSrc"
            :alt="`Google Cloud Certified - ${cert.name}`"
            class="h-36 w-36 md:h-40 md:w-40 object-contain transition-transform duration-300 group-hover:scale-105"
            width="263"
            height="263"
            format="webp"
            loading="lazy"
          />
        </a>
      </div>

      <p class="text-center text-muted mt-12">
        Also: DATAVERSITY Data Architecture Workshop: Practical Exercises in Architecture Design
        (October 2021, 24 CEUs).
        <a
          href="/pdfs/dataversity-data-architecture-workshop-2021.pdf"
          class="text-primary underline"
          target="_blank"
          rel="noopener"
          >View certificate (PDF)</a
        >
      </p>

      <div class="text-center mt-16">
        <a
          href="https://www.credly.com/users/nikolay-tonev.a60be769/badges#credly"
          class="bg-primary text-white font-bold py-3 px-8 rounded-lg hover:bg-opacity-90 transition-all"
          target="_blank"
          rel="noopener noreferrer"
          >View All on Credly</a
        >
      </div>
    </div>
  </section>
</template>
