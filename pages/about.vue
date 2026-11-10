<script setup lang="ts">
useSeoMeta({
  title: 'About | ivmanto.com',
  description:
    'Meet Nikolay Tonev, Cloud Data Architect & AI Consultant. 20+ years in enterprise data, GCP specialist, DAMA-aligned governance expert.',
  ogTitle: 'About | ivmanto.com',
  ogDescription:
    'Meet the expert behind IVMANTO — specializing in GCP data architecture, AI solutions, and cloud governance.',
  twitterTitle: 'About | ivmanto.com',
  twitterDescription:
    'Meet the expert behind IVMANTO — specializing in GCP data architecture, AI solutions, and cloud governance.',
})

useHead({
  script: [
    {
      id: 'person-schema',
      type: 'application/ld+json',
      innerHTML: JSON.stringify(
        {
          '@context': 'https://schema.org',
          '@type': 'Person',
          '@id': 'https://ivmanto.com/about#person',
          name: 'Nikolay Tonev',
          url: 'https://ivmanto.com/about',
          image: 'https://ivmanto.com/nt.jpg',
          jobTitle: 'Cloud Data Architect & AI Consultant',
          worksFor: { '@id': 'https://ivmanto.com/#organization' },
          sameAs: [
            'https://linkedin.com/in/nikolaytonev',
            'https://github.com/dasiyes',
            'https://g.dev/ivmanto-nikolaytonev',
          ],
          knowsAbout: [
            'Google Cloud Platform',
            'Data Architecture',
            'AI & Machine Learning',
            'Data Governance',
            'DAMA-DMBOK',
          ],
        },
        null,
        2,
      ),
    },
  ],
})
</script>

<template>
  <div>
    <section class="py-20 md:py-28 bg-light-gray">
      <div class="container mx-auto px-6 grid md:grid-cols-2 gap-12 items-center">
        <div>
          <NuxtImg
            src="/nt.jpg"
            alt="Nikolay Tonev - Cloud Data Architect and AI Consultant"
            class="rounded-xl shadow-lg w-full max-w-md mx-auto"
            width="448"
            height="448"
            format="webp"
            loading="lazy"
          />
        </div>
        <div>
          <h1 class="text-3xl md:text-4xl font-bold text-dark-slate">
            A Partner in Your Data Journey
          </h1>
          <p class="text-lg text-muted mt-6">
            I am a dedicated Cloud Data Architect with extensive experience on the Google Cloud
            Platform. My passion is helping businesses unlock the true value of their data. I focus
            on building solutions that are not just technologically sound, but also align perfectly
            with your business goals.
          </p>
          <p class="text-lg text-muted mt-4">
            I'm currently on the path to becoming a DAMA Certified Data Management Professional,
            ensuring that my work adheres to the highest industry standards for data governance and
            quality.
          </p>
          <div class="mt-8 flex flex-wrap gap-4">
            <a
              href="https://linkedin.com/in/nikolaytonev"
              class="btn-profile-link"
              target="_blank"
              rel="noopener noreferrer"
              >LinkedIn</a
            ><a
              href="https://github.com/dasiyes"
              class="btn-profile-link"
              target="_blank"
              rel="noopener noreferrer"
              >GitHub</a
            ><a
              href="https://www.credly.com/users/nikolay-tonev.a60be769/badges#credly"
              class="btn-profile-link"
              target="_blank"
              rel="noopener noreferrer"
              >Credly Profile</a
            ><a
              href="https://cdmp.info/"
              class="btn-profile-link"
              target="_blank"
              rel="noopener noreferrer"
              >DAMA Professionals</a
            >
          </div>
        </div>
      </div>
    </section>

    <AboutApproachSection />

    <AboutExpertiseSection />

    <AboutCertificatesSection />

    <!-- Featured Insights — internal cross-links for SEO -->
    <section class="py-16 md:py-20 bg-surface">
      <div class="container mx-auto px-6">
        <div class="text-center max-w-3xl mx-auto mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-dark-slate">Featured Insights</h2>
          <p class="text-lg text-muted mt-4">
            Explore my latest thinking on data strategy, cloud architecture, and AI.
          </p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          <NuxtLink
            to="/blog/the-shift-to-agentic-ai-and-autonomous-workflows"
            class="block p-6 bg-light-gray rounded-xl hover:shadow-lg transition-shadow"
          >
            <h3 class="text-lg font-bold text-dark-slate">The Shift to Agentic AI</h3>
            <p class="text-sm text-muted mt-2">
              How autonomous workflows are reshaping enterprise AI strategy.
            </p>
            <span class="mt-3 inline-block text-primary font-semibold text-sm">Read Article &rarr;</span>
          </NuxtLink>
          <NuxtLink
            to="/blog/DataMeshGovernance"
            class="block p-6 bg-light-gray rounded-xl hover:shadow-lg transition-shadow"
          >
            <h3 class="text-lg font-bold text-dark-slate">Data Mesh Governance</h3>
            <p class="text-sm text-muted mt-2">
              How to enforce policies without becoming a bottleneck.
            </p>
            <span class="mt-3 inline-block text-primary font-semibold text-sm">Read Article &rarr;</span>
          </NuxtLink>
          <NuxtLink
            to="/assessment"
            class="block p-6 bg-primary bg-opacity-5 rounded-xl border-2 border-primary border-opacity-20 hover:shadow-lg transition-shadow"
          >
            <h3 class="text-lg font-bold text-primary">Data Maturity Assessment</h3>
            <p class="text-sm text-muted mt-2">
              Take our free interactive assessment to evaluate your organization's data readiness.
            </p>
            <span class="mt-3 inline-block text-primary font-semibold text-sm">Start Assessment &rarr;</span>
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Call to Action Section -->

    <section class="bg-primary text-white">
      <div class="container mx-auto px-6 py-20 text-center">
        <h2 class="text-3xl md:text-4xl font-bold">Ready to Unlock Your Data's Potential?</h2>
        <p class="text-lg mt-4 max-w-2xl mx-auto">
          Let's discuss how we can build a data strategy that drives growth and innovation for your
          business.
        </p>
        <NuxtLink
          to="/#contact"
          class="mt-8 inline-block bg-surface text-primary font-bold py-3 px-8 rounded-lg hover:bg-opacity-90 transition-all"
          >Get In Touch</NuxtLink
        >
      </div>
    </section>
  </div>
</template>
