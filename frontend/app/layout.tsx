import type { Metadata } from "next";
import { ReactQueryProvider } from "@/lib/providers";
import "./globals.css";

export const metadata: Metadata = {
  title: "DAWAI - Assessment Platform",
  description: "Multi-subject assessment platform",
  manifest: "/manifest.json",
  themeColor: "#2563eb",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="manifest" href="/manifest.json" />
        <meta name="theme-color" content="#2563eb" />
        <meta name="mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      </head>
      <body>
        <ReactQueryProvider>{children}</ReactQueryProvider>
        <script
          dangerouslySetInnerHTML={{
            __html: `
              if ('serviceWorker' in navigator) {
                window.addEventListener('load', () => {
                  navigator.serviceWorker.register('/service-worker.js').catch(() => {});
                });
              }
            `,
          }}
        />
      </body>
    </html>
  );
}
