<template>
  <section class="py-20 md:py-28 bg-white">
    <div class="container mx-auto px-6">
      <div class="text-center max-w-3xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-bold text-dark-slate">Core Expertise</h2>
        <p class="text-lg text-gray-600 mt-4">
          I leverage a modern data stack to design and implement robust, scalable, and efficient
          data solutions that drive business value. My approach is rooted in a commitment to simple,
          maintainable designs, seamless team collaboration, and building with security and
          compliance by design.
        </p>
      </div>

      <div class="mt-16 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <!-- Data Management -->
        <div class="bg-light-gray p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow flex flex-col">
          <h3 class="text-xl font-bold text-primary mb-4">Data Management &amp; Governance</h3>
          <p class="text-gray-600 flex-grow">
            I provide comprehensive data management consulting to unlock the full potential of your
            data assets. My expertise covers the entire data lifecycle, from data governance,
            quality, and master data management to ensuring your data is a secure and compliant
            strategic asset for your organization.
          </p>
          <NuxtLink
            to="/services"
            class="mt-4 inline-block text-primary font-semibold hover:underline"
          >
            Explore this service &rarr;
          </NuxtLink>
        </div>

        <!-- Data Architecture -->
        <div class="bg-light-gray p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow flex flex-col">
          <h3 class="text-xl font-bold text-primary mb-4">Data Architecture</h3>
          <p class="text-gray-600 flex-grow">
            I design and build end-to-end data architectures tailored to your unique business goals.
            I guide you through the entire process, from initial strategy and solution design to
            full-scale implementation and production deployment, creating the foundational blueprint
            for your data-driven future.
          </p>
          <NuxtLink
            to="/services"
            class="mt-4 inline-block text-primary font-semibold hover:underline"
          >
            Explore this service &rarr;
          </NuxtLink>
        </div>

        <!-- Data Products -->
        <div class="bg-light-gray p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow flex flex-col">
          <h3 class="text-xl font-bold text-primary mb-4">Modern Data Products</h3>
          <p class="text-gray-600 flex-grow">
            I build and implement powerful, end-to-end data products that directly deliver business
            value. From design to development, I engineer bespoke solutions that integrate
            seamlessly with modern data platforms, including Data Warehouses, Data Lakes, and
            Lakehouses.
          </p>
          <NuxtLink
            to="/services"
            class="mt-4 inline-block text-primary font-semibold hover:underline"
          >
            Explore this service &rarr;
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>
