import { Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

export const MovementBarChart = ({ data }: { data: any[] }) => {
  if (!data || data.length === 0) return <p>Aucun mouvement récent.</p>;

  const chartData = {
    labels: data.map(d => d.date),
    datasets: [
      {
        label: 'Entrées',
        data: data.map(d => d.in_qty),
        backgroundColor: 'rgba(75, 192, 192, 0.5)',
      },
      {
        label: 'Sorties',
        data: data.map(d => d.out_qty),
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
      },
    ],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: { position: 'top' as const },
      title: { display: true, text: 'Mouvements de Stock (30 jours)' },
    },
    scales: {
        y: {
            beginAtZero: true
        }
    }
  };

  return <Bar options={options} data={chartData} />;
};
