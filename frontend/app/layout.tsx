import type { Metadata } from "next";
import { ReactQueryProvider } from "@/lib/providers";
import "./globals.css";

export const metadata: Metadata = {
  title: "DAWAI - Assessment Platform",
  description: "Multi-subject assessment platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ReactQueryProvider>{children}</ReactQueryProvider>
      </body>
    </html>
  );
}
