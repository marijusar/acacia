export const TicketCardPriority = () => {
  return Array(4)
    .fill(null)
    .map(() => (
      <div className="w-3 h-3 mr-2 rounded-full bg-(--destructive)" />
    ));
};
