import { footerSections, footerDescription, footerNote, footerPolicies } from "../data/footerLinks";


export const Footer: React.FC = () => {
  return (
    <footer className="bg-[#222222] text-white py-8 px-2 md:px-8 lg:px-12 font-[Quicksand] ">
      <div className="px-4">
        {/* Logo and Description */}
        <div className="mb-12">
          <div className="mb-6">
            <img
              src="/icons/LogoWhite.svg"
              alt="ALX Logo"
              className="w-[54px] h-[30px]"
            />
          </div>
          <p className="text-[#CACACA] text-[11px] font-medium md:text-[15px] leading-relaxed max-w-2xl">
            {footerDescription}
          </p>
        </div>

        {/* Footer Sections Grid */}
        <div className="grid grid-cols-2 md:grid-cols-3 gap-5 mb-12 pb-12 border-b border-[#FFFFFF17] text-[#CACACA] px-4">
          {footerSections.map((section) => (
            <div key={section.title}>
              <h3 className="text-[16px] md:text-[18px] font-semibold mb-6 text-white">
                {section.title}
              </h3>
              <ul className="space-y-2">
                {section.links.map((link) => (
                  <li key={link.label}>
                    <a
                      href={link.href}
                      className="text-[#CACACA] hover:text-white transition-colors text-[12px] md:text-[16px]"
                    >
                      {link.label}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <div className="lg:flex lg:justify-between  ">
          {/* Cancellation Note */}
          <div className=" text-center px-2 ">
            <p className=" text-[12px] md:text-[14px] text-[#CACACA]">
              {footerNote}
              <a
                href="#"
                className="text-[#34967C] hover:text-[#2a7a65] transition-colors ml-1"
              >
                here
              </a>
            </p>
          </div>

          {/* Policies */}
          <div className="flex flex-row justify-center items-center gap-6 md:gap-8 m-8 md:m-0">
            {footerPolicies.map((policy) => (
              <a
                key={policy.label}
                href={policy.href}
                className=" hover:text-white transition-colors text-[12px] md:text-[14px] text-[#CACACA]"
              >
                {policy.label}
              </a>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;