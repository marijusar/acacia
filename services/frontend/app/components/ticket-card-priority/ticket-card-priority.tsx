export const TicketCardPriority = () => {
  return Array(4)
    .fill(null)
    .map((_, i) => (
      <div key={i} className="w-3 h-3 mr-2 rounded-full bg-(--destructive)" />
    ));
};
